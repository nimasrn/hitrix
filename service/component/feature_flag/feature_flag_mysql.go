package featureflag

import (
	"errors"
	"fmt"

	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type IFeatureFlag interface {
	GetName() string
	ScriptsSingleInstance() []app.IScript
	ScriptsMultiInstance() []app.IScript
}

type serviceFeatureFlag struct {
	featureFlags       map[string]IFeatureFlag
	errorLoggerService errorlogger.ErrorLogger
}

func NewFeatureFlagService(errorLoggerService errorlogger.ErrorLogger) ServiceFeatureFlagInterface {
	featureFlags := make(map[string]IFeatureFlag)

	return &serviceFeatureFlag{
		featureFlags:       featureFlags,
		errorLoggerService: errorLoggerService,
	}
}

func (s *serviceFeatureFlag) IsActive(ormService fluxaorm.Context, name string) bool {
	if name == "" {
		panic("name cannot be empty")
	}

	//TODO Krasi ORM: fix it
	featureFlagEntity, found := fluxaorm.GetByUniqueIndex[entity.FeatureFlagEntity](ormService, "Name", name)
	if !found {
		return false
	}

	return featureFlagEntity.Enabled && featureFlagEntity.Registered
}

func (s *serviceFeatureFlag) FailIfIsNotActive(ormService fluxaorm.Context, name string) error {
	isActive := s.IsActive(ormService, name)
	if !isActive {
		return fmt.Errorf("feature (%s) is not active", name)
	}

	return nil
}

func (s *serviceFeatureFlag) Enable(ormService fluxaorm.Context, name string) error {
	if name == "" {
		panic("name cannot be empty")
	}

	//TODO Krasi ORM: fix it
	featureFlagEntity, found := fluxaorm.GetByUniqueIndex[entity.FeatureFlagEntity](ormService, "Name", name)
	if !found {
		return errors.New("feature cannot be found")
	}

	ormService.EditEntity(featureFlagEntity)

	featureFlagEntity.Enabled = true

	ormService.Flush()

	return nil
}

func (s *serviceFeatureFlag) Disable(ormService fluxaorm.Context, name string) error {
	if name == "" {
		panic("name cannot be empty")
	}

	//TODO Krasi ORM: fix it
	featureFlagEntity, found := fluxaorm.GetByUniqueIndex[entity.FeatureFlagEntity](ormService, "Name", name)
	if !found {
		return errors.New("feature cannot be found")
	}

	ormService.EditEntity(featureFlagEntity)

	featureFlagEntity.Enabled = false

	ormService.Flush()

	return nil
}

func (s *serviceFeatureFlag) getAllActive(ormService fluxaorm.Context, _ *fluxaorm.Pager) []IFeatureFlag {
	//TODO Krasi ORM: use pager fix it
	featureFlagEntitiesIterator := fluxaorm.GetByIndex[entity.FeatureFlagEntity](ormService, "Registered_Enabled", true, true)

	activeFeatureFlags := make([]IFeatureFlag, 0)

	for _, featureFlagEntity := range featureFlagEntitiesIterator.All() {
		if _, ok := s.featureFlags[featureFlagEntity.Name]; !ok {
			s.errorLoggerService.LogError("feature flag " + featureFlagEntity.Name + " is not registered")

			continue
		}

		activeFeatureFlags = append(activeFeatureFlags, s.featureFlags[featureFlagEntity.Name])
	}

	return activeFeatureFlags
}

func (s *serviceFeatureFlag) GetScriptsSingleInstance(ormService fluxaorm.Context) []app.IScript {
	activeFeatureFlags := s.getAllActive(ormService, fluxaorm.NewPager(1, 1000))

	allScripts := make([]app.IScript, 0)
	for _, featureFlag := range activeFeatureFlags {
		allScripts = append(allScripts, featureFlag.ScriptsSingleInstance()...)
	}

	return allScripts
}

func (s *serviceFeatureFlag) GetScriptsMultiInstance(ormService fluxaorm.Context) []app.IScript {
	activeFeatureFlags := s.getAllActive(ormService, fluxaorm.NewPager(1, 1000))

	allScripts := make([]app.IScript, 0)
	for _, featureFlag := range activeFeatureFlags {
		allScripts = append(allScripts, featureFlag.ScriptsMultiInstance()...)
	}

	return allScripts
}

func (s *serviceFeatureFlag) Register(featureFlags ...IFeatureFlag) {
	s.featureFlags = make(map[string]IFeatureFlag)

	for _, featureFlag := range featureFlags {
		if _, has := s.featureFlags[featureFlag.GetName()]; has {
			panic("feature flag with name '" + featureFlag.GetName() + "' already exists")
		}

		s.featureFlags[featureFlag.GetName()] = featureFlag
	}
}

func (s *serviceFeatureFlag) Sync(ormService fluxaorm.Context, clockService clock.IClock) {
	var (
		featureFlagEntities []*entity.FeatureFlagEntity
		lastID              uint64
	)

	for {
		pager := fluxaorm.NewPager(1, 1000)
		featureFlagEntitiesIterator := fluxaorm.Search[entity.FeatureFlagEntity](ormService, fluxaorm.NewWhere("ID > ? ORDER BY ID ASC", lastID), pager)

		rows := featureFlagEntitiesIterator.All()
		if len(rows) == 0 {
			break
		}

		lastID = rows[len(rows)-1].ID
		featureFlagEntities = append(featureFlagEntities, rows...)

		if len(rows) < pager.PageSize {
			break
		}
	}

	dbFeatureFlags := make(map[string]*entity.FeatureFlagEntity)

	for _, featureFlagEntity := range featureFlagEntities {
		if featureFlagEntity != nil {
			dbFeatureFlags[featureFlagEntity.Name] = featureFlagEntity
		} else {
			s.errorLoggerService.LogError("feature flag is nil")
		}
	}

	for _, registeredFeatureFlag := range s.featureFlags {
		if _, ok := dbFeatureFlags[registeredFeatureFlag.GetName()]; !ok {
			ormService.NewEntity(entity.FeatureFlagEntity{
				Name:       registeredFeatureFlag.GetName(),
				Registered: true,
				Enabled:    false,
				UpdatedAt:  nil,
				CreatedAt:  clockService.Now(),
			})

			err := ormService.FlushWithCheck()
			if err != nil {
				if duplicateKeyError, ok := err.(*fluxaorm.DuplicateKeyError); ok {
					if duplicateKeyError.Index != "Name" {
						panic(err)
					}
				} else {
					panic(err)
				}
			}
		} else if !dbFeatureFlags[registeredFeatureFlag.GetName()].Registered {
			featureFlagEntity := dbFeatureFlags[registeredFeatureFlag.GetName()]

			ormService.NewEntity(featureFlagEntity)

			featureFlagEntity.Registered = true
			featureFlagEntity.UpdatedAt = clockService.NowPointer()
		}
	}

	for name, dbFeatureFlag := range dbFeatureFlags {
		if _, ok := s.featureFlags[name]; !ok {
			ormService.EditEntity(dbFeatureFlag)

			dbFeatureFlag.Registered = false
			dbFeatureFlag.UpdatedAt = clockService.NowPointer()
		}
	}

	ormService.Flush()
}
