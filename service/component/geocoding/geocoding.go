package geocoding

import (
	"context"
	//nolint //G501: Blocklisted import crypto/md5: weak cryptographic primitive, but just fine for caching
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/latolukasz/fluxaorm"
	"googlemaps.github.io/maps"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/clock"
)

type IGeocoding interface {
	SnapToRoad(ctx context.Context, dto *maps.SnapToRoadRequest) (*maps.SnapToRoadResponse, error)
	Geocode(ctx context.Context, ormService fluxaorm.Context, address string, language string) (*Address, error)
	ReverseGeocode(ctx context.Context, ormService fluxaorm.Context, latLng *LatLng, language string) (*Address, error)
	CutCoordinates(float float64, precision int) (float64, error)
}

type Address struct {
	Found                    bool
	FromCache                bool
	AdministrativeAreaLevel1 string
	CityName                 string
	Address                  string
	Language                 string
	Location                 *LatLng
}

type LatLng struct {
	Lat float64
	Lng float64
}

type Geocoding struct {
	useCaching      bool
	cacheTTLMinDays int
	cacheTTLMaxDays int
	clock           clock.IClock
	provider        Provider
}

func NewGeocoding(
	useCaching bool,
	cacheTTLMinDays int,
	cacheTTLMaxDays int,
	clock clock.IClock,
	provider Provider,
) IGeocoding {
	return &Geocoding{
		useCaching:      useCaching,
		cacheTTLMinDays: cacheTTLMinDays,
		cacheTTLMaxDays: cacheTTLMaxDays,
		clock:           clock,
		provider:        provider,
	}
}

func (g *Geocoding) SnapToRoad(ctx context.Context, dto *maps.SnapToRoadRequest) (*maps.SnapToRoadResponse, error) {
	return g.provider.SnapToRoad(ctx, dto)
}

func (g *Geocoding) Geocode(ctx context.Context, ormService fluxaorm.Context, address string, language string) (*Address, error) {
	address = strings.TrimSpace(address)

	if g.useCaching {
		geocodingEntity, found := fluxaorm.GetByUniqueIndex[entity.GeocodingCacheEntity](
			ormService,
			"AddressHash_Language",
			g.getAddressHash(address),
			language,
		)

		if found {
			return &Address{
				Found:                    true,
				FromCache:                true,
				AdministrativeAreaLevel1: geocodingEntity.AdministrativeAreaLevel1,
				CityName:                 geocodingEntity.CityName,
				Address:                  address,
				Language:                 geocodingEntity.Language,
				Location: &LatLng{
					Lat: geocodingEntity.Lat,
					Lng: geocodingEntity.Lng,
				},
			}, nil
		}
	}

	geocodedAddress, providerRawResponse, err := g.provider.Geocode(ctx, address, language)
	if err != nil {
		return nil, err
	}

	if g.useCaching && geocodedAddress.Found {
		now := g.clock.Now()

		ormService.NewEntity(&entity.GeocodingCacheEntity{
			Lat:                      geocodedAddress.Location.Lat,
			Lng:                      geocodedAddress.Location.Lng,
			AdministrativeAreaLevel1: geocodedAddress.AdministrativeAreaLevel1,
			CityName:                 geocodedAddress.CityName,
			Address:                  address,
			AddressHash:              g.getAddressHash(address),
			Language:                 language,
			Provider:                 g.provider.GetName(),
			RawResponse:              providerRawResponse,
			ExpiresAt:                now.Add(time.Duration(g.getCacheTTL(g.cacheTTLMinDays, g.cacheTTLMaxDays)) * time.Hour * 24),
			CreatedAt:                now,
		})

		ormService.Flush()
	}

	return geocodedAddress, nil
}

func (g *Geocoding) ReverseGeocode(ctx context.Context, ormService fluxaorm.Context, latLng *LatLng, language string) (*Address, error) {
	cacheLat := latLng.Lat
	cacheLng := latLng.Lng

	if g.useCaching {
		var err error

		cacheLat, err = g.CutCoordinates(cacheLat, 5)
		if err != nil {
			return nil, err
		}

		cacheLng, err = g.CutCoordinates(cacheLng, 5)
		if err != nil {
			return nil, err
		}

		reverseGeocodingCacheEntity, found := fluxaorm.GetByUniqueIndex[entity.GeocodingReverseCacheEntity](
			ormService,
			"Lat_Lng_Language",
			cacheLng,
			cacheLng,
			language,
		)
		if found {
			return &Address{
				Found:                    true,
				FromCache:                true,
				AdministrativeAreaLevel1: reverseGeocodingCacheEntity.AdministrativeAreaLevel1,
				CityName:                 reverseGeocodingCacheEntity.CityName,
				Address:                  reverseGeocodingCacheEntity.Address,
				Language:                 reverseGeocodingCacheEntity.Language,
				Location: &LatLng{
					Lat: reverseGeocodingCacheEntity.Lat,
					Lng: reverseGeocodingCacheEntity.Lng,
				},
			}, nil
		}
	}

	geocodedAddress, providerRawResponse, err := g.provider.ReverseGeocode(ctx, latLng, language)
	if err != nil {
		return nil, err
	}

	if g.useCaching && geocodedAddress.Found {
		now := g.clock.Now()

		ormService.NewEntity(&entity.GeocodingReverseCacheEntity{
			Lat:                      cacheLat,
			Lng:                      cacheLng,
			AdministrativeAreaLevel1: geocodedAddress.AdministrativeAreaLevel1,
			CityName:                 geocodedAddress.CityName,
			Address:                  strings.TrimSpace(geocodedAddress.Address),
			Language:                 language,
			Provider:                 g.provider.GetName(),
			RawResponse:              providerRawResponse,
			ExpiresAt:                now.Add(time.Duration(g.getCacheTTL(g.cacheTTLMinDays, g.cacheTTLMaxDays)) * time.Hour * 24),
			CreatedAt:                now,
		})

		ormService.Flush()
	}

	return geocodedAddress, nil
}

func (g *Geocoding) CutCoordinates(float float64, precision int) (float64, error) {
	asString := fmt.Sprintf("%.8f", float)
	asStringParts := strings.Split(asString, ".")

	return strconv.ParseFloat(asString[0:len(asStringParts[0])+1+precision], 64)
}

func (g *Geocoding) getCacheTTL(cacheTTLMinDays, cacheTTLMaxDays int) int {
	intVal := big.NewInt(int64(cacheTTLMaxDays - cacheTTLMinDays))

	randomInt, err := rand.Int(rand.Reader, intVal)
	if err != nil {
		panic(err)
	}

	return int(randomInt.Int64() + int64(cacheTTLMinDays))
}

func (g *Geocoding) getAddressHash(address string) string {
	//nolint // G401: Use of weak cryptographic primitive , but just fine for caching
	hash := md5.Sum([]byte(address))

	return hex.EncodeToString(hash[:])
}
