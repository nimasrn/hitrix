package translation

import (
	"context"

	"github.com/latolukasz/fluxaorm"

	listDto "github.com/coretrix/hitrix/pkg/dto/list"
	"github.com/coretrix/hitrix/pkg/dto/translation"
	"github.com/coretrix/hitrix/pkg/entity"
	crudView "github.com/coretrix/hitrix/pkg/view/crud"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/crud"
)

const (
	pageSizeMin = 10
	pageSizeMax = 100
)

func columns() []*crud.Column {
	return []*crud.Column{
		{
			Key:                      "ID",
			FilterType:               crud.InputTypeNumber,
			Label:                    "ID",
			Searchable:               false,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:        "Status",
			FilterType: crud.SelectTypeStringString,
			Label:      "Status",
			Searchable: true,
			Sortable:   false,
			Visible:    true,
			DataStringKeyStringValue: []*crud.StringKeyStringValue{
				{Key: entity.TranslationStatusNew.String(), Label: entity.TranslationStatusNew.String()},
				{Key: entity.TranslationStatusTranslated.String(), Label: entity.TranslationStatusTranslated.String()},
			},
		},
		{
			Key:        "Lang",
			FilterType: crud.InputTypeString,
			Label:      "Lang",
			Searchable: true,
			Sortable:   false,
			Visible:    true,
		},
		{
			Key:        "Key",
			FilterType: crud.InputTypeString,
			Label:      "Key",
			Searchable: true,
			Sortable:   false,
			Visible:    true,
		},
		{
			Key:        "Text",
			Label:      "Text",
			Searchable: false,
			Sortable:   false,
			Visible:    true,
		},
	}
}

func List(ctx context.Context, userListRequest listDto.RequestDTOList) (*translation.ResponseDTOList, error) {
	request, err := crudView.ValidateListRequest(userListRequest, pageSizeMin, pageSizeMax)
	if err != nil {
		return nil, err
	}

	cols := columns()
	crudService := service.DI().Crud()

	searchParams := crudService.ExtractListParams(cols, request)
	where := crudService.GenerateListMysqlQuery(searchParams)

	if len(searchParams.Sort) == 0 {
		where.Append("ORDER BY ID DESC")
	}

	ormService := service.DI().OrmForContext(ctx)

	entityIterator, total := fluxaorm.SearchWithCount[entity.TranslationTextEntity](
		ormService,
		where,
		fluxaorm.NewPager(searchParams.Page, searchParams.PageSize),
	)
	rows := make([]*translation.ListRow, entityIterator.Len())

	for i, translationTextEntity := range entityIterator.All() {
		rows[i] = &translation.ListRow{
			ID:     translationTextEntity.ID,
			Status: translationTextEntity.Status,
			Lang:   translationTextEntity.Lang,
			Key:    translationTextEntity.Key,
			Text:   translationTextEntity.Text,
		}
	}

	return &translation.ResponseDTOList{
		Rows:    rows,
		Total:   int(total),
		Columns: cols,
	}, nil
}
