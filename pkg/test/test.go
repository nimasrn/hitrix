package test

import (
	"fmt"
	"math/rand"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/fluxaorm"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

var createTableExecuted = false
var parallelTestID string

type Environment struct {
	t                *testing.T
	Hitrix           *hitrix.Hitrix
	GinEngine        *gin.Engine
	Cxt              *gin.Context
	ResponseRecorder *httptest.ResponseRecorder
}

func CreateContext(
	t *testing.T,
	projectName string,
	defaultServices []*service.DefinitionGlobal,
	mockGlobalServices []*service.DefinitionGlobal,
	redisPools *app.RedisPools,
) *Environment {
	var deferFunc func()

	err := os.Setenv("TZ", "UTC")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("APP_MODE", app.ModeTest)
	if err != nil {
		t.Fatal(err)
	}

	testSpringInstance, deferFunc := hitrix.New(projectName, "").
		SetParallelTestID(getParallelID()).
		RegisterDIGlobalService(append(defaultServices, mockGlobalServices...)...).RegisterRedisPools(
		redisPools,
	).Build()
	defer deferFunc()

	ormService := service.DI().Orm()

	executeAlters(ormService)

	return &Environment{t: t, Hitrix: testSpringInstance}
}

func CreateAPIContext(
	t *testing.T,
	projectName string,
	ginInitHandler hitrix.GinInitHandler,
	defaultGlobalServices []*service.DefinitionGlobal,
	defaultRequestServices []*service.DefinitionRequest,
	mockGlobalServices []*service.DefinitionGlobal,
	mockRequestServices []*service.DefinitionRequest,
	redisPools *app.RedisPools,
) *Environment {
	var deferFunc func()

	err := os.Setenv("TZ", "UTC")
	if err != nil {
		t.Fatal(err)
	}

	err = os.Setenv("APP_MODE", app.ModeTest)
	if err != nil {
		t.Fatal(err)
	}

	testSpringInstance, deferFunc := hitrix.New(projectName, "").
		SetParallelTestID(getParallelID()).
		RegisterDIGlobalService(append(defaultGlobalServices, mockGlobalServices...)...).
		RegisterDIRequestService(append(defaultRequestServices, mockRequestServices...)...).
		RegisterRedisPools(
			redisPools,
		).Build()

	defer deferFunc()

	ginTestInstance := hitrix.InitGin(ginInitHandler)

	ormService := service.DI().Orm()

	executeAlters(ormService)

	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)

	return &Environment{t: t, Hitrix: testSpringInstance, GinEngine: ginTestInstance, Cxt: c, ResponseRecorder: resp}
}

func executeAlters(ormService fluxaorm.Context) {
	if !createTableExecuted {
		createTableExecuted = true

		alters := fluxaorm.GetAlters(ormService)

		for pool, db := range ormService.Engine().Registry().DBPools() {
			dbAlters := ""

			dropTables(ormService, db)

			for _, alter := range alters {
				if alter.Pool == pool {
					dbAlters += alter.SQL
				}
			}

			if dbAlters != "" {
				db.Exec(ormService, dbAlters)
			}
		}
	} else {
		for _, db := range ormService.Engine().Registry().DBPools() {
			truncateTables(ormService, db)
		}
	}

	if os.Getenv("PARALLEL_TESTS") == "" || os.Getenv("PARALLEL_TESTS") == "false" {
		for _, localCache := range ormService.Engine().Registry().LocalCachePools() {
			localCache.Clear(ormService)
		}

		for _, redisCache := range ormService.Engine().Registry().RedisPools() {
			redisCache.FlushAll(ormService)
		}
	}

	altersSearch := fluxaorm.GetRedisSearchAlters(ormService)
	for _, alter := range altersSearch {
		alter.Exec(ormService)
	}
}

func getRandomString() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 10)

	//nolint //G404: Use of weak random number generator (math/rand instead of crypto/rand)
	rand.Read(b)

	return fmt.Sprintf("%x%d", b, os.Getpid())[:5]
}

func getParallelID() string {
	if os.Getenv("PARALLEL_TESTS") == "" || os.Getenv("PARALLEL_TESTS") == "false" {
		return "1"
	} else if parallelTestID != "" {
		return parallelTestID
	}

	parallelTestID = getRandomString()

	return parallelTestID
}

func dropTables(ormService fluxaorm.Context, db fluxaorm.DB) {
	var query string

	rows, deferF := db.Query(ormService,
		"SELECT CONCAT('DROP TABLE IF EXISTS ',table_schema,'.',table_name,';') AS query "+
			"FROM information_schema.tables WHERE table_schema IN ('"+db.GetConfig().GetDatabaseName()+"')",
	)

	defer deferF()

	if rows != nil {
		var queries string

		for rows.Next() {
			rows.Scan(&query)
			queries += query
		}

		_, def := db.Query(ormService, "SET FOREIGN_KEY_CHECKS=0;"+queries+"SET FOREIGN_KEY_CHECKS=1")

		defer def()
	}
}

// TODO Krasi ORM: delete -> truncate
func truncateTables(ormService fluxaorm.Context, db fluxaorm.DB) {
	var query string

	rows, deferF := db.Query(ormService,
		"SELECT CONCAT('delete from  ',table_schema,'.',table_name,';' , 'ALTER TABLE ', table_schema,'.',table_name , ' AUTO_INCREMENT = 1;') AS query "+
			"FROM information_schema.tables WHERE table_schema IN ('"+db.GetConfig().GetDatabaseName()+"');",
	)

	defer deferF()

	if rows != nil {
		var queries string

		for rows.Next() {
			rows.Scan(&query)
			queries += query
		}

		_, def := db.Query(ormService, "SET FOREIGN_KEY_CHECKS=0;"+queries+"SET FOREIGN_KEY_CHECKS=1")
		defer def()
	}
}
