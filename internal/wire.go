//go:build wireinject
// +build wireinject

package internal

//go:generate go run github.com/google/wire/cmd/wire
import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/hhr0815hhr/gint/internal/controller"
	"github.com/hhr0815hhr/gint/internal/database/model"
	"github.com/hhr0815hhr/gint/internal/database/mysql"
	"github.com/hhr0815hhr/gint/internal/logic"
	"github.com/hhr0815hhr/gint/internal/server/http"
)

type AppInfo struct {
	Engine    *gin.Engine
	TestLogic *logic.TestLogic
	Data      map[string]interface{}
}

func ProvideApp(
	engine *gin.Engine,
	testLogic *logic.TestLogic,
) *AppInfo {
	return &AppInfo{
		Engine:    engine,
		TestLogic: testLogic,
		Data:      map[string]interface{}{},
	}
}

func InitApp() *AppInfo {
	wire.Build(DbSet, RepoSet, LogicSet, RouteSet, ProvideRoutes, http.NewHTTPServer, ProvideApp)
	return nil
}

var DbSet = wire.NewSet(mysql.ProvideDB)

var RepoSet = wire.NewSet(
	model.NewTestRepo,
)
var LogicSet = wire.NewSet(
	logic.NewTestLogic,
)

var RouteSet = wire.NewSet(
	controller.NewTestController,
)

func ProvideRoutes(
	test *controller.TestController,
) *http.HTTPRoutes {
	return &http.HTTPRoutes{
		Routers: []controller.Router{
			test,
		},
	}
}
