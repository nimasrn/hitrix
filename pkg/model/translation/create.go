package translation

import (
	"context"
	"errors"

	"github.com/coretrix/hitrix/pkg/dto/translation"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/latolukasz/fluxaorm"
)

func Create(ctx context.Context, request *translation.RequestCreateTranslation) (*translation.ResponseTranslation, error) {
	ormService := service.DI().OrmForContext(ctx)

	newTranslationEntity := entity.TranslationTextEntity{
		Lang:   request.Lang.String(),
		Key:    request.Key.String(),
		Status: entity.TranslationStatusTranslated.String(),
		Text:   request.Text,
	}

	fluxaorm.NewEntityFromSource(ormService, newTranslationEntity)

	//TODO Krasi ORM: check for error
	err := ormService.Flush()
	if err != nil {
		return nil, errors.New("text with this lang and key already exists")
		//return nil, errors.HandleFlushWithCheckError(
		//	err,
		//	errors.HandleCustomErrors(map[string]string{"Lang": "translation text with this lang and key already exists"}),
		//)
	}

	return &translation.ResponseTranslation{
		ID:   newTranslationEntity.ID,
		Lang: newTranslationEntity.Lang,
		Key:  newTranslationEntity.Key,
		Text: newTranslationEntity.Text,
	}, nil
}
