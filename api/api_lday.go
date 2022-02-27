package api

import (
    "fmt"
    "container/list"

    "github.com/kataras/iris/v12"

    "github.com/c-doge/stock.go/db"
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
    Type   string    `url:"type"`
    Adjust string    `url:"adjust"`
    Code   string    `url:"code"`
    From   string    `url:"from"`
    To     string    `url:"to"`
    Head   bool      `url:"head"`
}

func apiV1PutLday(ctx iris.Context) {
    var req PutLdayRequest;
    err := ctx.ReadProtobuf(&req)
    if err != nil {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": err.Error()})
        return;
    }
    if req.Data == nil || len(req.Data) == 0 {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": "KData is empty"})
        return;
    }
    var m = gostk.NewKDataMap()
    for _, k := range req.Data {
        code, kdata := convertProtoBufToKData(k)
        if gostk.IsStockCodeValid(code) && kdata.Time.After(utils.EarlyDate) {
//            logger.Infof("[API.Lday] ------%s: %v", code, kdata);
            m.Insert(code, kdata)
        }
    }

    var total_stock int = 0;
    var total_kdata int = 0;
    m.ForEach(func(code string, l *list.List) {
        ll := make([]*gostk.KData, 0, l.Len())
        for e := l.Front(); e != nil; e = e.Next() {
            kk := e.Value.(*gostk.KData)
            ll = append(ll, kk)
        }
        if len(ll) > 0 {
            err = db.PutLday(code, ll)
            if err != nil {
                logger.Warnf("[API.Lday] db.PutLday fail, error: %v", err);
            } else {
                total_kdata += len(ll);
                total_stock += 1;
            }
        }
    })
    msg := fmt.Sprintf("%d stock and %d kdata updated", total_stock, total_kdata)
    logger.Infof("[API.Lday] PutLday %s", msg)
    
    ctx.JSON(iris.Map{"status": StatusOK, "message": msg})
}

func apiV1GetLday(ctx iris.Context) {
    var req GetLdayRequest;

    err := ctx.ReadQuery(&req)
    if err != nil && !iris.IsErrPath(err) {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusServerError, "message": err.Error()})
        return;
    }
    if req.Type != "csv" {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": "Accept Type Unknown"})
        return
    }
    if !gostk.IsStockCodeValid(req.Code) {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": gostk.ErrorStockCodeUnknown.Error()})
        return
    }
    from, e1 := utils.ParseTime("20060102", req.From)
    to, e2 := utils.ParseTime("20060102", req.To)
    if e1 != nil || e2 != nil {
        msg := ""
        if e1 != nil {
            msg = fmt.Sprintf("Parse %s fail, Error: %s", req.From, e1);
        } else {
            msg = fmt.Sprintf("Parse %s fail, Error: %s", req.To, e2);
        }
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": msg})
        return
    }

    if from.After(to) {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": "From time is later than To time"})
        return
    }
    
    l, err := db.GetLday(req.Code, from, to);
    if err != nil {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusServerError, "message": err.Error()})
        return
    }

    ctx.StatusCode(iris.StatusOK)
    ctx.Header("Transfer-Encoding", "chunked")
    if req.Head {
        ctx.WriteString("Date,Open,Close,High,Low,Volume,Turnover\n")
    }
    for _, k := range l {
        ctx.Writef("%s,%.4f,%.4f,%.4f,%.4f,%.4f,%.4f\n",
                                        k.Time.Format("2006-01-02"),
                                        k.Open,
                                        k.Close,
                                        k.High,
                                        k.Low,
                                        k.Volume,
                                        k.Turnover);
    }
}
