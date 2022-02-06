package api

import (
	"time"
	"github.com/kataras/iris/v12"
)

func apiHealth(ctx iris.Context) {
	
	ts := time.Now().Unix()
	ctx.JSON(iris.Map{"status": iris.StatusOK, "timestamp": ts})
}