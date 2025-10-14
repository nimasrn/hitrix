package ddos

import (
	"strconv"
	"time"

	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/service/component/app"
)

type IDDOS interface {
	ProtectManyAttempts(appService *app.App, ormService fluxaorm.Context, protectCriterion string, maxAttempts int, ttl int) bool
}

type DDOS struct {
}

func (t *DDOS) ProtectManyAttempts(appService *app.App, ormService fluxaorm.Context, protectCriterion string, maxAttempts int, ttl int) bool {
	redis := ormService.Engine().Redis(appService.RedisPools.Cache)
	attempts, has := redis.Get(ormService, "ddos_"+protectCriterion)
	count := 0

	if len(attempts) > 0 {
		var err error

		count, err = strconv.Atoi(attempts)
		if err != nil {
			panic(err)
		}
	}

	if has && count >= maxAttempts {
		return false
	}

	redis.Set(ormService, "ddos_"+protectCriterion, count+1, time.Duration(ttl)*time.Second)

	return true
}
