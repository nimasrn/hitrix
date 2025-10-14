package scripts

import (
	"context"
	"log"
	"os"

	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

type DBSeedScript struct {
	SeedsPerProject map[string][]Seed
}

func (script *DBSeedScript) Run(_ context.Context, _ app.IExit, ormService fluxaorm.Context) {
	appService := service.DI().App()
	Seeder(script.SeedsPerProject, ormService, appService)
}

func (script *DBSeedScript) Infinity() bool {
	return true
}

func (script *DBSeedScript) Unique() bool {
	return true
}

func (script *DBSeedScript) Description() string {
	return "Seed Database"
}

type Seed interface {
	Execute(fluxaorm.Context)
	Environments() []string
	Name() string
}

func Seeder(seedsPerProject map[string][]Seed, ormService fluxaorm.Context, appService *app.App) {
	for project, seeds := range seedsPerProject {
		if project != os.Getenv("PROJECT_NAME") {
			continue
		}

		for _, seed := range seeds {
			supportCurrentEnv := false

			for _, env := range seed.Environments() {
				if env == appService.Mode {
					supportCurrentEnv = true

					break
				}
			}

			if !supportCurrentEnv {
				continue
			}

			whereStmt := fluxaorm.NewWhere("`Name` = ?", seed.Name())

			seederEntity, found := fluxaorm.SearchOne[entity.SeederEntity](ormService, whereStmt)

			if found {
				continue
			}

			seed.Execute(ormService)

			seederEntity.Name = seed.Name()
			seederEntity.CreatedAt = service.DI().Clock().Now()

			fluxaorm.EditEntity(ormService, seederEntity)
			ormService.Flush()

			log.Println("Seeder " + seed.Name() + " has been executed")
		}
	}
}
