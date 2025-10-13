package translation

import (
	"context"
	"errors"
	"fmt"

	"github.com/coretrix/hitrix/pkg/dto/translation"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/latolukasz/fluxaorm"
)

func Update(ctx context.Context, request *translation.RequestUpdateTranslation, id uint64) (*translation.ResponseTranslation, error) {
	ormService := service.DI().OrmForContext(ctx)

	translationTextEntity, found := fluxaorm.GetByID[entity.TranslationTextEntity](ormService, id)
	if !found {
		return nil, fmt.Errorf("translation text with ID %v not found", id)
	}

	ormService.EditEntity(translationTextEntity)

	translationTextEntity.Lang = request.Lang.String()
	translationTextEntity.Key = request.Key.String()
	translationTextEntity.Text = request.Text
	translationTextEntity.Status = entity.TranslationStatusTranslated.String()

	//TODO Krasi ORM: check for error
	err := ormService.Flush()
	if err != nil {
		return nil, errors.New("translation text with this lang and key already exists")
		//return nil, errors.HandleFlushWithCheckError(
		//	err,
		//	errors.HandleCustomErrors(map[string]string{"Lang": "translation text with this lang and key already exists"}),
		//)
	}

	return &translation.ResponseTranslation{
		ID:   translationTextEntity.ID,
		Lang: translationTextEntity.Lang,
		Key:  translationTextEntity.Key,
		Text: translationTextEntity.Text,
	}, nil
}
