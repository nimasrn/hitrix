package scripts

import (
	"context"
	"time"

	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

type ClearExpiredGeocodingCache struct {
}

func (script *ClearExpiredGeocodingCache) Run(_ context.Context, _ app.IExit, ormService fluxaorm.Context) {
	now := service.DI().Clock().Now()

	fiveAM := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location())

	if now.After(fiveAM) {
		//TODO sleep few hours
		return
	}

	where := fluxaorm.NewWhere("ExpiresAt < ?", now)

	entityIterator := fluxaorm.Search[entity.GeocodingCacheEntity](ormService, where, fluxaorm.NewPager(1, 10000))

	for _, geocodingEntity := range entityIterator.All() {
		fluxaorm.DeleteEntity(ormService, geocodingEntity)
	}

	ormService.Flush()

	where = fluxaorm.NewWhere("ExpiresAt < ?", now)

	entityIterator1 := fluxaorm.Search[entity.GeocodingReverseCacheEntity](ormService, where, fluxaorm.NewPager(1, 10000))

	for _, reverseGeocodingCacheEntity := range entityIterator1.All() {
		fluxaorm.DeleteEntity(ormService, reverseGeocodingCacheEntity)
	}

	ormService.Flush()
}

func (script *ClearExpiredGeocodingCache) Interval() time.Duration {
	return time.Hour * 3
}

func (script *ClearExpiredGeocodingCache) Unique() bool {
	return true
}

func (script *ClearExpiredGeocodingCache) Description() string {
	return "clear expired geocoding cache"
}
