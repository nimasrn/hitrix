package mocks

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/latolukasz/fluxaorm"
	"github.com/stretchr/testify/mock"
	"googlemaps.github.io/maps"

	"github.com/coretrix/hitrix/service/component/geocoding"
)

type FakeGeocoding struct {
	mock.Mock
}

func (f *FakeGeocoding) Geocode(_ context.Context, _ fluxaorm.Context, address string, language string) (*geocoding.Address, error) {
	args := f.Called(address, language)

	return args.Get(0).(*geocoding.Address), args.Error(1)
}

func (f *FakeGeocoding) SnapToRoad(ctx context.Context, dto *maps.SnapToRoadRequest) (*maps.SnapToRoadResponse, error) {
	args := f.Called(ctx, dto)

	return args.Get(0).(*maps.SnapToRoadResponse), args.Error(1)
}

func (f *FakeGeocoding) ReverseGeocode(
	_ context.Context,
	_ fluxaorm.Context,
	latLng *geocoding.LatLng,
	language string,
) (*geocoding.Address, error) {
	args := f.Called(latLng, language)

	return args.Get(0).(*geocoding.Address), args.Error(1)
}

func (f *FakeGeocoding) CutCoordinates(float float64, precision int) (float64, error) {
	asString := fmt.Sprintf("%.8f", float)
	asStringParts := strings.Split(asString, ".")

	return strconv.ParseFloat(asString[0:len(asStringParts[0])+1+precision], 64)
}
