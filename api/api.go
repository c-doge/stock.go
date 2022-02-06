package api

import (
	"github.com/kataras/iris/v12"

	"github.com/c-doge/stock.go/base/setting"
    "github.com/c-doge/stock.go/base/logger"
)

var webApp *iris.Application = nil

func Init(run bool) error {
    logger.Info("API Start")

    webApp = iris.New()
	api := webApp.Party("/api")

    api.Get("/health",        apiHealth);
    api.Post("/health",        apiHealth);

   // api.Get("/lday",     apiGetLday)
    api.Put("/lday",     apiPutLday)
    // app.Delete("/lday",  apiDelLDay)

    if run {
    	err := webApp.Listen(setting.Get().Web.Addr)
    	return err;
    }
    return nil;
}