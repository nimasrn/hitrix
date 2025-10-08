package featureflag

import (
	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
)

type ServiceFeatureFlagInterface interface {
	IsActive(ormService fluxaorm.Context, name string) bool
	FailIfIsNotActive(ormService fluxaorm.Context, name string) error
	Enable(ormService fluxaorm.Context, name string) error
	Disable(ormService fluxaorm.Context, name string) error
	GetScriptsSingleInstance(ormService fluxaorm.Context) []app.IScript
	GetScriptsMultiInstance(ormService fluxaorm.Context) []app.IScript
	Register(featureFlags ...IFeatureFlag)
	Sync(ormService fluxaorm.Context, clockService clock.IClock)
}
