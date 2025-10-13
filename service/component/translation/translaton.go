package translation

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/pkg/helper"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type ITranslationService interface {
	GetText(ormService fluxaorm.Context, lang entity.TranslationTextLang, key entity.TranslationTextKey) string
	GetTextWithVars(
		ormService fluxaorm.Context,
		lang entity.TranslationTextLang,
		key entity.TranslationTextKey,
		variables map[string]interface{},
	) string
}

type translationService struct {
	errorLoggerService errorlogger.ErrorLogger
}

func NewTranslationService(errorLoggerService errorlogger.ErrorLogger) ITranslationService {
	return &translationService{errorLoggerService}
}

func (u *translationService) GetText(ormService fluxaorm.Context, lang entity.TranslationTextLang, key entity.TranslationTextKey) string {
	translationTextEntity, found := fluxaorm.GetByUniqueIndex[entity.TranslationTextEntity](
		ormService,
		"Lang_Key",
		lang.String(),
		key.String(),
	)
	if !found {
		newTranslationTextEntity := entity.TranslationTextEntity{
			Lang:   lang.String(),
			Key:    key.String(),
			Status: entity.TranslationStatusNew.String(),
		}

		ormService.NewEntity(newTranslationTextEntity)

		ormService.Flush()

		return key.String()
	}

	if translationTextEntity.Status == entity.TranslationStatusNew.String() {
		return key.String()
	}

	re := regexp.MustCompile(`\[\[(.*?)\]\]`)
	subMatchAll := re.FindAllString(translationTextEntity.Text, -1)

	if len(subMatchAll) > 0 {
		u.errorLoggerService.LogError(
			fmt.Sprintf(
				"not assigned vars (%s) for translation key `%s`",
				strings.Join(subMatchAll, ", "),
				key.String(),
			),
		)
	}

	return translationTextEntity.Text
}

func (u *translationService) GetTextWithVars(
	ormService fluxaorm.Context,
	lang entity.TranslationTextLang,
	key entity.TranslationTextKey,
	variables map[string]interface{},
) string {
	keys := make([]string, 0, len(variables))

	for k := range variables {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	translationTextEntity, found := fluxaorm.GetByUniqueIndex[entity.TranslationTextEntity](
		ormService,
		"Lang_Key",
		lang.String(),
		key.String(),
	)
	if !found {
		newTranslationTextEntity := entity.TranslationTextEntity{
			Lang:   lang.String(),
			Key:    key.String(),
			Status: entity.TranslationStatusNew.String(),
			Vars:   keys,
		}

		ormService.NewEntity(newTranslationTextEntity)

		ormService.Flush()

		return key.String()
	}

	if !helper.EqualString(translationTextEntity.Vars, keys) {
		ormService.EditEntity(translationTextEntity)

		translationTextEntity.Vars = keys

		ormService.Flush()
	}

	if translationTextEntity.Status == entity.TranslationStatusNew.String() {
		return key.String()
	}

	text := translationTextEntity.Text

	for paramName, value := range variables {
		text = strings.Replace(text, fmt.Sprintf("[[%s]]", paramName), fmt.Sprintf("%v", value), -1)
	}

	re := regexp.MustCompile(`\[\[(.*?)\]\]`)
	subMatchAll := re.FindAllString(text, -1)

	if len(subMatchAll) > 0 {
		u.errorLoggerService.LogError(
			fmt.Sprintf(
				"not assigned vars (%s) for translation key `%s`",
				strings.Join(subMatchAll, ", "),
				key.String(),
			),
		)
	}

	return text
}
