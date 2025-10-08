package entity

import (
	"time"
)

type GeocodingReverseCacheEntity struct {
	ID                       uint64  `orm:"table=geocoding_reverse_cache;localCache;redisCache"`
	Lat                      float64 `orm:"decimal=8,5;required;unique=Lat_Lng_Language:1;cached"`
	Lng                      float64 `orm:"decimal=8,5;required;unique=Lat_Lng_Language:2;cached"`
	AdministrativeAreaLevel1 string  `orm:"required"`
	CityName                 string  `orm:"required"`
	Address                  string
	Language                 string `orm:"required;enum=entity.LanguageValueAll;unique=Lat_Lng_Language:3;cached"`
	Provider                 string
	RawResponse              interface{}
	ExpiresAt                time.Time `orm:"time=true;index=ExpiresAt"`
	CreatedAt                time.Time `orm:"time=true"`
}
