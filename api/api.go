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
	api := webApp.Party("/api/v1")

    api.Get("/health",         apiV1Health);
    api.Post("/health",        apiV1Health);

    api.Get("/lday",     apiV1GetLday)
    api.Put("/lday",     apiV1PutLday)

    api.Get("/vol",     apiV1GetVolumeList)
    api.Put("/vol",     apiV1PutVolumeList) 

    api.Get("/xdr",     apiV1GetXDRList)
    api.Put("/xdr",     apiV1PutXDRList)    
    if run {
    	err := webApp.Listen(setting.Get().Web.Addr)
    	return err;
    }
    return nil;
}