package entity

import (
	"time"
)

type Color string

var Colors = struct {
	Red    Color
	Blue   Color
	Yellow Color
}{
	Red:    "red",
	Blue:   "blue",
	Yellow: "yellow",
}

func (c Color) EnumValues() any {
	return Colors
}

type GeocodingCacheEntity struct {
	ID                       uint64  `orm:"table=geocoding_cache;localCache;redisCache"`
	Lat                      float64 `orm:"decimal=8,5"`
	Lng                      float64 `orm:"decimal=8,5"`
	AdministrativeAreaLevel1 string  `orm:"required"`
	CityName                 string  `orm:"required"`
	Address                  string  `orm:"required"`
	AddressHash              string  `orm:"length=32;required;unique=AddressHash_Language:1;cached"`
	Language                 string  `orm:"required;enum=entity.LanguageValueAll;unique=AddressHash_Language:2;cached"`
	Provider                 string
	RawResponse              interface{}
	ExpiresAt                time.Time `orm:"time=true;index=ExpiresAt"`
	CreatedAt                time.Time `orm:"time=true"`
}
