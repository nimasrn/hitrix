package setting

import (
	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

type ServiceSettingInterface interface {
	Get(ormService fluxaorm.Context, key string) (*entity.SettingsEntity, bool)
	GetString(ormService fluxaorm.Context, key string) (string, bool)
	GetInt(ormService fluxaorm.Context, key string) (int, bool)
	GetUint(ormService fluxaorm.Context, key string) (uint, bool)
	GetInt64(ormService fluxaorm.Context, key string) (int64, bool)
	GetUint64(ormService fluxaorm.Context, key string) (uint64, bool)
	GetFloat64(ormService fluxaorm.Context, key string) (float64, bool)
	GetBool(ormService fluxaorm.Context, key string) (bool, bool)
}
