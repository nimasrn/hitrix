package translation

import (
	"context"
	"fmt"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/latolukasz/fluxaorm"
)

func Delete(ctx context.Context, id uint64) error {
	ormService := service.DI().OrmForContext(ctx)

	translationTextEntity, found := fluxaorm.GetByID[entity.TranslationTextEntity](ormService, id)
	if !found {
		return fmt.Errorf("translation text with ID %v not found", id)
	}

	fluxaorm.DeleteEntity(ormService, translationTextEntity)

	ormService.Flush()

	return nil
}
