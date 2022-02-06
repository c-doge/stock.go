package api

import (
//  "time"
    "github.com/kataras/iris/v12"

    "github.com/c-doge/stock.go/gostk"
    "github.com/c-doge/stock.go/base/utils"
    "github.com/c-doge/stock.go/base/logger"
)

func convertProtoBufToKData(k *KData) (string, *gostk.KData) {
    kk := &gostk.KData {
        Time:       utils.DecimalNumToDateTime(k.Time),
        Open:       k.GetOpen(),
        Close:      k.GetClose(),
        High:       k.GetHigh(),
        Low:        k.GetLow(),
        Volume:     k.GetVolume(),
        Turnover:   k.GetTurnover(),
    }
    code := k.GetCode()
    return code, kk;
}

type GetLdayRequest struct {
    Type  string    `url:"type"`
    Code  string    `url:"code"`
    From  string    `url:"from"`
    To    string    `url:"to"`
}

func apiPutLday(ctx iris.Context) {
    var req PutLdayRequest;
    err := ctx.ReadProtobuf(&req)
    if err != nil {
        ctx.StopWithError(iris.StatusBadRequest, err)
        return
    }
    if req.Data == nil || len(req.Data) == 0 {
        ctx.StopWithError(iris.StatusBadRequest, err)
        return
    }
    var m = gostk.NewKDataMap()
    for _, k := range req.Data {
        code, kdata := convertProtoBufToKData(k)
        logger.Infof("%s: %v", code, kdata);
        m.Insert(code, kdata)
    }
    
    ctx.JSON(iris.Map{"status": iris.StatusOK})

}

func apiGetLday(ctx iris.Context) {
    var req GetLdayRequest;

    err := ctx.ReadQuery(&req)
    if err != nil && !iris.IsErrPath(err) {
        ctx.StatusCode(iris.StatusInternalServerError)
        ctx.WriteString(err.Error())
    }
//    ctx.JSON(iris.Map{"status": iris.StatusOK})
    ctx.Writef("MyType: %#v", req)
}
