package api

import (
    "fmt"
    "time"
    "github.com/kataras/iris/v12"

    "github.com/c-doge/stock.go/db"
    "github.com/c-doge/stock.go/gostk"
    "github.com/c-doge/stock.go/base/utils"
//    "github.com/c-doge/stock.go/base/logger"
)

type GetVolRequest struct {
    Type   string    `url:"type"`
    Code   string    `url:"code"`
    From   string    `url:"from"`
    To     string    `url:"to"`
    Head   bool      `url:"head"`
}

type GetXDRRequest struct {
    Type   string    `url:"type"`
    Code   string    `url:"code"`
    From   string    `url:"from"`
    To     string    `url:"to"`
    Head   bool      `url:"head"`
}

func apiV1GetVolumeList(ctx iris.Context) {
	var req GetVolRequest;
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
    var from time.Time = utils.EarlyDate
    var to   time.Time = utils.FutureDate
    if len(req.From) != 0 {
    	from, err = utils.ParseTime("20060102", req.From)
    	if err != nil {
    		msg := fmt.Sprintf("Parse %s fail, Error: %s", req.From, err);
    		ctx.StatusCode(iris.StatusOK)
        	ctx.JSON(iris.Map{"status": StatusBadRequest, "message": msg})
        	return
    	}
    }
    if len(req.To) != 0 {
    	to, err = utils.ParseTime("20060102", req.To)
    	if err != nil {
    		msg := fmt.Sprintf("Parse %s fail, Error: %s", req.To, err);
    		ctx.StatusCode(iris.StatusOK)
        	ctx.JSON(iris.Map{"status": StatusBadRequest, "message": msg})
        	return
    	}
    }
    if from.After(to) {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": "From time is later than To time"})
        return
    }
    l, err := db.GetVolumeList(req.Code, from, to)
    if err != nil {
    	ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusServerError, "message": err.Error()})
        return
    }
    if req.Head {
    	ctx.WriteString("Date,PreTradable,PostTradable,PreTotal,PostTotal\n")
    }
   	for _, k := range l {
        ctx.Writef("%s,%.6f,%.6f,%.6f,%.6f\n",
                                        k.Date.Format("2006-01-02"),
                                        k.PreTradable,
                                        k.PostTradable,
                                        k.PreTotal,
                                        k.PostTotal);
    }
}

func apiV1PutVolumeList(ctx iris.Context) {
	var req PutVolRequest;
    err := ctx.ReadProtobuf(&req)
    if err != nil {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": err.Error()})
        return;
    }
    if req.Data == nil || len(req.Data) == 0 {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": "VData is empty"})
        return;
    }
    if !gostk.IsStockCodeValid(req.Code) {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": "Stock Code Unknown"})
        return;
    }
    l := make([]*gostk.VData, len(req.Data))
    for i := 0; i < len(req.Data); i++ {
    	l[i] = &gostk.VData {
    		Date:         utils.DecimalNumToDateTime(req.Data[i].Date),
    		PreTradable:  req.Data[i].GetPreTradable(),
    		PreTotal:     req.Data[i].GetPreTotal(),
			PostTradable: req.Data[i].GetPostTradable(),
    		PostTotal:    req.Data[i].GetPostTotal(),
    	}
    }
    err = db.PutVolumeList(req.Code, l)
    if err != nil {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusServerError, "message": err.Error()})
        return
    }
    ctx.StatusCode(iris.StatusOK)
    ctx.JSON(iris.Map{"status": StatusOK})
}


func apiV1GetXDRList(ctx iris.Context) {
	var req GetXDRRequest
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
    var from time.Time = utils.EarlyDate
    var to   time.Time = utils.FutureDate
    if len(req.From) != 0 {
    	from, err = utils.ParseTime("20060102", req.From)
    	if err != nil {
    		msg := fmt.Sprintf("Parse %s fail, Error: %s", req.From, err);
    		ctx.StatusCode(iris.StatusOK)
        	ctx.JSON(iris.Map{"status": StatusBadRequest, "message": msg})
        	return
    	}
    }
    if len(req.To) != 0 {
    	to, err = utils.ParseTime("20060102", req.To)
    	if err != nil {
    		msg := fmt.Sprintf("Parse %s fail, Error: %s", req.To, err);
    		ctx.StatusCode(iris.StatusOK)
        	ctx.JSON(iris.Map{"status": StatusBadRequest, "message": msg})
        	return
    	}
    }
    if from.After(to) {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": "From time is later than To time"})
        return
    }
    l, err := db.GetXDRList(req.Code, from, to)
    if err != nil {
    	ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusServerError, "message": err.Error()})
        return
    }
    if req.Head {
    	ctx.WriteString("Date,Type,AllotPrice,AllotVolume,BonusPrice,BonusVolume\n")
    }
   	for _, k := range l {
        ctx.Writef("%s,%s,%.6f,%.6f,%.6f,%.6f\n",
                                        k.Date.Format("2006-01-02"),
                                        k.Type.String(),
                                        k.AllotPrice,
                                        k.AllotVolume,
                                        k.BonusPrice,
                                        k.BonusVolume);
    }
}

func apiV1PutXDRList(ctx iris.Context) {
    getType := func (t XDataType) gostk.XDataType {
        if t == XData_EXP {
            return gostk.XData_EXP
        }
        return gostk.XData_XDR
    }
	var req PutXDRRequest;
    err := ctx.ReadProtobuf(&req)
    if err != nil {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": err.Error()})
        return;
    }
    if req.Data == nil || len(req.Data) == 0 {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": "XData is empty"})
        return;
    }
    if !gostk.IsStockCodeValid(req.Code) {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusBadRequest, "message": "Stock Code Unknown"})
        return;
    }
    l := make([]*gostk.XData, len(req.Data))
    for i := 0; i < len(req.Data); i++ {
    	l[i] = &gostk.XData {
    		Date:        utils.DecimalNumToDateTime(req.Data[i].Date),
            Type:        getType(req.Data[i].GetType()),
    		AllotPrice:  req.Data[i].GetAllotPrice(),
    		AllotVolume: req.Data[i].GetAllotVolume(),
			BonusPrice:  req.Data[i].GetBonusPrice(),
    		BonusVolume: req.Data[i].GetBonusVolume(),
    	}
    }
    err = db.PutXDRList(req.Code, l)
    if err != nil {
        ctx.StatusCode(iris.StatusOK)
        ctx.JSON(iris.Map{"status": StatusServerError, "message": err.Error()})
        return
    }
    ctx.StatusCode(iris.StatusOK)
    ctx.JSON(iris.Map{"status": StatusOK, "message": "success"})
}
