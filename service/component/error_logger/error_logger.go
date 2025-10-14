package errorlogger

import (
	"bytes"
	//nolint //G501: Blocklisted import crypto/md5: weak cryptographic primitive
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http/httputil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/fluxaorm"
	slackgo "github.com/slack-go/slack"

	"github.com/coretrix/hitrix/service/component/app"
	requestlogger "github.com/coretrix/hitrix/service/component/request_logger"
	"github.com/coretrix/hitrix/service/component/sentry"
	"github.com/coretrix/hitrix/service/component/slack"
)

const GroupError = "error"
const GroupWarning = "warning"

type eventConfig struct {
	redisKey     string
	title        string
	anchor       string
	slackChannel string
}

type EventRow struct {
	ID      string
	Type    string
	File    string
	Line    int
	AppName string
	Request string
	Message string
	Stack   string
	Counter int
	Time    string
}

type event struct {
	File    string
	Line    int
	AppName string
	Request []byte
	Message string
	Stack   []byte
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type ErrorLogger interface {
	LogError(errData interface{})
	LogErrorWithRequest(c *gin.Context, errData interface{})
	LogWarning(errData interface{})
	LogWarningWithRequest(c *gin.Context, errData interface{})
	LogPanicWithRequest(c *gin.Context, errData interface{})
	GetErrors() []EventRow
	DeleteError(id string)
	DeleteAllErrors()
	GetWarnings() []EventRow
	DeleteWarning(id string)
	DeleteAllWarnings()
}

type RedisErrorLogger struct {
	redis          fluxaorm.RedisCache
	sentryService  sentry.ISentry
	slackService   slack.Slack
	appService     *app.App
	ormService     fluxaorm.Context
	requestBodyKey interface{}
}

func NewRedisErrorLogger(
	appService *app.App,
	ormService fluxaorm.Context,
	redis fluxaorm.RedisCache,
	slackService slack.Slack,
	sentryService sentry.ISentry,
	requestBodyKey interface{},
) ErrorLogger {
	return &RedisErrorLogger{
		appService:     appService,
		ormService:     ormService,
		redis:          redis,
		slackService:   slackService,
		sentryService:  sentryService,
		requestBodyKey: requestBodyKey,
	}
}

func (e *RedisErrorLogger) LogError(errData interface{}) {
	e.log(errData, 2, nil, false)
}

func (e *RedisErrorLogger) LogErrorWithRequest(c *gin.Context, errData interface{}) {
	e.log(errData, 2, c, false)
}

func (e *RedisErrorLogger) LogWarning(errData interface{}) {
	e.log(errData, 2, nil, true)
}

func (e *RedisErrorLogger) LogWarningWithRequest(c *gin.Context, errData interface{}) {
	e.log(errData, 2, c, true)
}

func (e *RedisErrorLogger) LogPanicWithRequest(c *gin.Context, errData interface{}) {
	e.log(errData, 4, c, false)
}

func (e *RedisErrorLogger) GetErrors() []EventRow {
	return e.get(GroupError)
}

func (e *RedisErrorLogger) DeleteError(id string) {
	e.redis.HDel(e.ormService, GroupError, id)
	e.redis.HDel(e.ormService, GroupError, id+":time")
	e.redis.HDel(e.ormService, GroupError, id+":counter")
}

func (e *RedisErrorLogger) DeleteAllErrors() {
	e.redis.Del(e.ormService, GroupError)
}

func (e *RedisErrorLogger) GetWarnings() []EventRow {
	return e.get(GroupWarning)
}

func (e *RedisErrorLogger) DeleteWarning(id string) {
	e.redis.HDel(e.ormService, GroupWarning, id)
	e.redis.HDel(e.ormService, GroupWarning, id+":time")
	e.redis.HDel(e.ormService, GroupWarning, id+":counter")
}

func (e *RedisErrorLogger) DeleteAllWarnings() {
	e.redis.Del(e.ormService, GroupWarning)
}

func (e *RedisErrorLogger) log(errData interface{}, callerSkip int, c *gin.Context, warning bool) {
	var msg string

	err, ok := errData.(error)
	if ok {
		msg = err.Error()
	} else {
		msg = errData.(string)
	}

	logger := log.New(os.Stderr, "\n\n\x1b[31m", log.LstdFlags)
	stackTrace := stack(0)
	logger.Printf("[Error]:\n%s\n%s%s", msg, stackTrace, "\033[0m")

	_, file, line, _ := runtime.Caller(callerSkip)

	//nolint //G401: Use of weak cryptographic primitive
	errorKeyBinary := md5.Sum([]byte(e.appService.Name + ":" + file + ":" + fmt.Sprint(line)))
	errorKey := hex.EncodeToString(errorKeyBinary[:])
	value := &event{
		File:    file,
		Line:    line,
		AppName: e.appService.Name,
		Message: msg,
		Stack:   stackTrace,
	}

	if c != nil {
		c.Request.Body = io.NopCloser(bytes.NewReader(c.Request.Context().Value(e.requestBodyKey).([]byte)))

		requestID, has := c.Get(requestlogger.ID)
		if has {
			value.Request = []byte("X-Request-ID: " + fmt.Sprint(requestID) + "\n\n")
		}

		binaryRequest, _ := httputil.DumpRequest(c.Request, true)
		if len(binaryRequest)*4 <= 64000 {
			value.Request = append(value.Request, binaryRequest...)
		} else {
			value.Request = append(value.Request, []byte("Partial BODY \n\n")...)
			value.Request = append(value.Request, binaryRequest[0:16000]...)
		}
	}

	marshalValue, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}

	config := e.getEventConfig(warning)

	e.redis.HSet(e.ormService, config.redisKey, errorKey, marshalValue)
	e.redis.HSet(e.ormService, config.redisKey, errorKey+":time", time.Now().Unix())
	counter := e.redis.HIncrBy(e.ormService, config.redisKey, errorKey+":counter", 1)

	logg := math.Log10(float64(counter))

	if (e.slackService != nil && !e.appService.IsInLocalMode() && !e.appService.IsInTestMode()) && logg == float64(int64(logg)) {
		_ = e.slackService.SendToChannel(
			"errors",
			config.slackChannel,
			value.Message,
			slackgo.MsgOptionAttachments(
				slackgo.Attachment{
					AuthorName: e.appService.Name,
					Title:      config.title,
					TitleLink:  e.slackService.GetDevPanelURL() + "#" + config.anchor + "-" + errorKey,
					Text:       "Counter: " + fmt.Sprint(counter) + " ENV: " + e.appService.Mode,
				},
			),
		)
	}

	if (e.sentryService != nil && !e.appService.IsInLocalMode() && !e.appService.IsInTestMode()) &&
		logg == float64(int64(logg)) {
		e.sentryService.CaptureException(fmt.Errorf("%s", value.Message))
	}
}

func (e *RedisErrorLogger) get(group string) []EventRow {
	eventsData := e.redis.HGetAll(e.ormService, group)

	eventsList := map[string]*EventRow{}

	for key, value := range eventsData {
		// TODO: fix this hack
		if len(value) == 0 {
			continue
		}

		splitKeys := strings.Split(key, ":")

		if _, ok := eventsList[splitKeys[0]]; !ok {
			eventsList[splitKeys[0]] = &EventRow{
				ID: splitKeys[0],
			}
		}

		if len(splitKeys) == 1 {
			eventData := &event{}

			err := json.Unmarshal([]byte(value), eventData)
			if err != nil {
				eventsList[splitKeys[0]].Type = group
				eventsList[splitKeys[0]].Request = ""
				eventsList[splitKeys[0]].Stack = "unmarshal panic"
				eventsList[splitKeys[0]].File = "error_logger.go"
				eventsList[splitKeys[0]].Message = err.Error()
				eventsList[splitKeys[0]].Line = 0
				eventsList[splitKeys[0]].AppName = "ErrorLogger"

				continue
			}

			eventsList[splitKeys[0]].Type = group
			eventsList[splitKeys[0]].Request = string(eventData.Request)
			eventsList[splitKeys[0]].Stack = string(eventData.Stack)
			eventsList[splitKeys[0]].File = eventData.File
			eventsList[splitKeys[0]].Message = eventData.Message
			eventsList[splitKeys[0]].Line = eventData.Line
			eventsList[splitKeys[0]].AppName = eventData.AppName
		} else if len(splitKeys) == 2 {
			switch splitKeys[1] {
			case "time":
				i, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					i = 0
				}

				eventsList[splitKeys[0]].Time = time.Unix(i, 0).String()
			case "counter":
				counter, err := strconv.Atoi(value)
				if err != nil {
					counter = 0
				}

				eventsList[splitKeys[0]].Counter = counter
			}
		}
	}

	list := make([]EventRow, len(eventsList))
	for _, value := range eventsList {
		list = append(list, *value)
	}

	return list
}

func (e *RedisErrorLogger) getEventConfig(warning bool) eventConfig {
	if warning {
		return eventConfig{
			redisKey:     GroupWarning,
			title:        "Warning Link",
			anchor:       "#warn",
			slackChannel: e.slackService.GetErrorChannel(),
		}
	}

	return eventConfig{
		redisKey:     GroupError,
		title:        "Error Link",
		anchor:       "#err",
		slackChannel: e.slackService.GetErrorChannel(),
	}
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer)

	var (
		lines    [][]byte
		lastFile string
	)

	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)

		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}

		_, _ = fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}

	return buf.Bytes()
}

func source(lines [][]byte, n int) []byte {
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}

	return bytes.TrimSpace(lines[n])
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}

	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}

	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}

	name = bytes.ReplaceAll(name, centerDot, dot)

	return name
}
