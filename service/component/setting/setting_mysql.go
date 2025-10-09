package setting

import (
	"strconv"
	"strings"

	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/pkg/entity"
)

type serviceSetting struct {
	cache map[string]*entity.SettingsEntity
}

func NewSettingService() ServiceSettingInterface {
	return &serviceSetting{cache: map[string]*entity.SettingsEntity{}}
}

func (s *serviceSetting) Get(ormService fluxaorm.Context, key string) (*entity.SettingsEntity, bool) {
	if cachedEntity, exists := s.cache[key]; exists {
		return cachedEntity, true
	}

	settingEntity, found := fluxaorm.GetByUniqueIndex[entity.SettingsEntity](ormService, "SettingsKey", key)
	if !found {
		return nil, false
	}

	if !settingEntity.Editable && !settingEntity.Deletable {
		s.cache[key] = settingEntity
	}

	return settingEntity, true
}

func (s *serviceSetting) GetString(ormService fluxaorm.Context, key string) (string, bool) {
	setting, found := s.Get(ormService, key)
	if found {
		return setting.Value, true
	}

	return "", false
}

func (s *serviceSetting) GetInt(ormService fluxaorm.Context, key string) (int, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return 0, false
	}

	i, err := strconv.ParseInt(setting.Value, 10, 64)
	if err != nil {
		return 0, false
	}

	return int(i), true
}

func (s *serviceSetting) GetUint(ormService fluxaorm.Context, key string) (uint, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return 0, false
	}

	i, err := strconv.ParseUint(setting.Value, 10, 64)
	if err != nil {
		return 0, false
	}

	return uint(i), true
}

func (s *serviceSetting) GetInt64(ormService fluxaorm.Context, key string) (int64, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return 0, false
	}

	i, err := strconv.ParseInt(setting.Value, 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func (s *serviceSetting) GetUint64(ormService fluxaorm.Context, key string) (uint64, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return 0, false
	}

	i, err := strconv.ParseUint(setting.Value, 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func (s *serviceSetting) GetFloat64(ormService fluxaorm.Context, key string) (float64, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return 0, false
	}

	i, err := strconv.ParseFloat(setting.Value, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func (s *serviceSetting) GetBool(ormService fluxaorm.Context, key string) (bool, bool) {
	setting, found := s.Get(ormService, key)
	if !found {
		return false, false
	}

	if strings.ToLower(setting.Value) == "false" {
		return false, true
	}

	return true, true
}
