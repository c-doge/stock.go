package api

import (
    "io"
    "fmt"
    "time"
    "math"
    "bytes"
    "strings"
    "strconv"
    "testing"
  //  "net/http"
    "math/rand"
    "encoding/csv"
    "encoding/json"
    htest "net/http/httptest"
    "github.com/golang/protobuf/proto"
    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/httptest"

    "github.com/c-doge/stock.go/gostk"
    "github.com/c-doge/stock.go/base/utils"
    "github.com/c-doge/stock.go/base/logger"
)

func parseFloat(v string) float64 {
    s, _ := strconv.ParseFloat(v, 64)
    return s
}

func equalKData(d1, d2 *KData) bool {
    if d1.GetTime() != d2.GetTime() {
        logger.Debugf("time: %d,%d", d1.GetTime(), d2.GetTime())
        return false
    } else if abs32(d1.GetOpen(), d2.GetOpen()) > 0.0001 {
        logger.Debugf("open: %f,%f", d1.GetOpen(), d2.GetOpen())
        return false
    } else if abs32(d1.GetClose(), d2.GetClose()) > 0.0001 {
        logger.Debugf("close: %f,%f", d1.GetClose(), d2.GetClose())
        return false
    } else if abs32(d1.GetHigh(), d2.GetHigh()) > 0.0001 {
        logger.Debugf("high: %f,%f", d1.GetHigh(), d2.GetHigh())
        return false
    } else if abs32(d1.GetLow(), d2.GetLow()) > 0.0001 {
        logger.Debugf("low: %f,%f", d1.GetLow(), d2.GetLow())
        return false
    } else if math.Abs(d1.GetVolume() - d2.GetVolume()) > 0.0001 {
        logger.Debugf("Volume, %f-%f", d1.GetVolume(), d2.GetVolume())
        return false
    } else if math.Abs(d1.GetTurnover() - d2.GetTurnover()) > 0.0001 {
        logger.Debugf("Turnover: %f-%f", d1.GetTurnover(), d2.GetTurnover())
        return false
    }
    return true;
}

func abs32(f1, f2 float32) float32 {
    var v float32 = f1 - f2;
    if v > 0 {
        return v;
    }
    return 0-v;
}

func newStkKData(t time.Time) *gostk.KData {
    kdata := &gostk.KData {
        Time:      t,
        Open:      rand.Float32() * 100,
    };
    if rand.Intn(10) > 5 {
        kdata.Close = kdata.Open + kdata.Open * 0.1 * rand.Float32();
        kdata.High = kdata.Close + kdata.Open * 0.1 * rand.Float32();
        kdata.Low = kdata.Open - kdata.Open * 0.1 * rand.Float32();
    } else {
        kdata.Close = kdata.Open - kdata.Open * 0.1 * rand.Float32();
        kdata.High = kdata.Open + kdata.Open * 0.1 * rand.Float32();
        kdata.Low = kdata.Close - kdata.Open * 0.1 * rand.Float32();
    }
    kdata.Volume = rand.Float64() * 100000;
    kdata.Turnover = kdata.Volume * float64(kdata.Close);
    return kdata;
}

func newOndDayPutLdayRequest(dateNum uint32, n int) *PutLdayRequest {
    dateTime := utils.DecimalNumToDateTime(dateNum)
    l_kdata := make([]*KData, 0, n)
    for i:= 0; i < n; i++ {
        code := fmt.Sprintf("sh%6d", 600000+i+1)
        k := newStkKData(dateTime)
        kdata := &KData {
            Code:     code,
            Time:     utils.DateTimeToDecimalNum(dateTime),
            Open:     k.Open,
            Close:    k.Close,
            High:     k.High,
            Low:      k.Low,
            Volume:   k.Volume,
            Turnover: k.Turnover,
        }
        l_kdata = append(l_kdata, kdata)
    }
    req := &PutLdayRequest {
        Data: l_kdata,
    }
    return req
}

func newOneStockPutLdayRequest(code string, dateNum uint32, n int) *PutLdayRequest {
    l := make([]*KData, 0, n)
    dateTime := utils.DecimalNumToDateTime(dateNum)
    logger.Infof("code: %s, dateNum: %d (%s), n: %d", code, dateNum, dateTime.Format("2006-01-02"), n)
    for i:= 0; i < n; i++ {
        k := newStkKData(dateTime)
        kdata := &KData {
            Code:     code,
            Time:     utils.DateTimeToDecimalNum(dateTime),
            Open:     k.Open,
            Close:    k.Close,
            High:     k.High,
            Low:      k.Low,
            Volume:   k.Volume,
            Turnover: k.Turnover,
        }
        l = append(l, kdata)
        dateTime = dateTime.AddDate(0, 0, 1)
    }
    req := &PutLdayRequest {
        Data: l,
    }
    return req
}

func putLday(req *PutLdayRequest) error {

    var m ResponseModel;
    reqBin, err := proto.Marshal(req);
    if err != nil {
        //logger.Warnf("protobuf marshal LDayRequest fail, Error: %v", err);
        return err;
    }
    w := httptest.NewRecorder()
    r := httptest.NewRequest("PUT", "/api/v1/lday", bytes.NewReader(reqBin))
    httptest.Do(w, r, apiV1PutLday)

    err = json.Unmarshal(w.Body.Bytes(), &m);
    if err != nil {
        //logger.Warnf("putLday, Unmarshal response fail, Error: %s", err);
        return err;
    }
    if m.Status != StatusOK {
        return fmt.Errorf("Response Not Success, statusCode %d", m.Status)
    }
    return nil;
}

func getLday(code string, from, to uint32) (*csv.Reader, error) {
    path := fmt.Sprintf("/api/v1/lday?type=csv&code=%s&from=%d&to=%d", code, from, to)
    r := httptest.NewRequest("GET", path, nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiV1GetLday)

    if w.Code != iris.StatusOK {
        return nil, fmt.Errorf("Response Status Code: %d", w.Code)
    }

    reader := csv.NewReader(w.Body)
    return reader, nil
}

func checkCsvWithPutMsg(code string, r *csv.Reader, ll []*KData) error {
    index := 0;
    for {
        record, err := r.Read();
        if err == io.EOF {
            break
        } else if err != nil {
            return fmt.Errorf("Read CSV Fail, Error: %s", err);
        }
        if index >= len(ll) {
            return fmt.Errorf("the kdata got is more than put, index: %d", index);
        }
        t, err := utils.ParseTime("2006-01-02", record[0])
        if err != nil {
            return fmt.Errorf("parse Date fail, error: %v", err);
        }
        k := &KData {
            Code:     code,
            Time:     utils.DateTimeToDecimalNum(t),
            Open:     float32(parseFloat(record[1])),
            Close:    float32(parseFloat(record[2])),
            High:     float32(parseFloat(record[3])),
            Low:      float32(parseFloat(record[4])),
            Volume:   parseFloat(record[5]),
            Turnover: parseFloat(record[6]),
        }
        if !equalKData(k, ll[index]) {
            return fmt.Errorf("csv[%d] != putReq.Data[%d], Index:%v", k.Time, ll[index].Time, index);
        }
        index += 1;
    }
    if index != len(ll) {
        return fmt.Errorf("lines of csv(%d) not equal to the size of putReq.Data", index);
    }
    return nil
}

func checkResponseError(w *htest.ResponseRecorder, status int, errMsg string) error {

    var m ResponseModel
    if w.Code != iris.StatusOK {
        return fmt.Errorf("check statusCode fail: %d", w.Code);
    }
    err := json.Unmarshal(w.Body.Bytes(), &m);
    if err != nil {
        return fmt.Errorf("Unmarshal response json fail, Error: %v", err);
    }
    if m.Status != status {
        return fmt.Errorf("check reponse status fail, status: %d", m.Status);
    }
    if len(errMsg) > 0 && !strings.Contains(m.Message, errMsg) {
        return fmt.Errorf("check response error message fail, %s", m.Message)
    }
    return nil
}

func Test_GetLday(t *testing.T) {

    // bad type
    caseName := "BadType"
    r1 := httptest.NewRequest("GET", "/api/v1/lday?type=csx&code=sh000001&from=20120101&to=20120131", nil)
    r1.Header.Add("Accept", "text/csv")
    w1 := httptest.NewRecorder()
    httptest.Do(w1, r1, apiV1GetLday)
//    t.Logf("%s ==> \nResponse: %d,\nBody:%s", caseName, w1.Code, w1.Body.String());
    err := checkResponseError(w1, StatusBadRequest, "Accept Type Unknown")
    if err != nil {
        t.Errorf("TestGetLday.%s test fail, Error: %v", caseName, err);
        return
    }

    // bad code
    caseName = "BadCode"
    r2 := httptest.NewRequest("GET", "/api/v1/lday?type=csv&code=sa000001&from=20120101&to=20120131", nil)
    r2.Header.Add("Accept", "text/csv")
    w2 := httptest.NewRecorder()
    httptest.Do(w2, r2, apiV1GetLday)
//    t.Logf("%s ==> \nResponse: %d,\nBody:\n%s", caseName, w2.Code, w2.Body.String());
    err = checkResponseError(w2, StatusBadRequest, "Stock Code Unknown")
    if err != nil {
        t.Errorf("TestGetLday.%s test fail, Error: %v", caseName, err);
        return
    }

    // bad time
    caseName = "BadFromTime"
    r3 := httptest.NewRequest("GET", "/api/v1/lday?type=csv&code=sh000001&from=120120101&to=20120131", nil)
    r3.Header.Add("Accept", "text/csv")
    w3 := httptest.NewRecorder()
    httptest.Do(w3, r3, apiV1GetLday)
//    t.Logf("%s: \nResponse: %d,\nBody:\n%s", caseName, w3.Code, w3.Body.String());
    err = checkResponseError(w3, StatusBadRequest, "Parse 120120101 fail")
    if err != nil {
        t.Errorf("TestGetLday.%s test fail, Error: %v", caseName, err);
        return
    }

    caseName = "BadToTime"
    r4 := httptest.NewRequest("GET", "/api/v1/lday?type=csv&code=sh000001&from=20120101&to=201201311", nil)
    r4.Header.Add("Accept", "text/csv")
    w4 := httptest.NewRecorder()
    httptest.Do(w4, r4, apiV1GetLday)
//    t.Logf("%s ==> \nResponse: %d,\nBody:\n%s", caseName, w4.Code, w4.Body.String());
    err = checkResponseError(w4, StatusBadRequest, "Parse 201201311 fail")
    if err != nil {
        t.Errorf("TestGetLday.%s test fail, Error: %v", caseName, err);
        return
    }

    caseName = "FromLaterThanTo"
    r5 := httptest.NewRequest("GET", "/api/v1/lday?type=csv&code=sh000001&from=20120101&to=20110131", nil)
    r5.Header.Add("Accept", "text/csv")
    w5 := httptest.NewRecorder()
    httptest.Do(w5, r5, apiV1GetLday)
//    t.Logf("%s ==> \nResponse: %d,\nBody:\n%s", caseName, w5.Code, w5.Body.String());
    err = checkResponseError(w5, StatusBadRequest, "From time is later than To time")
    if err != nil {
        t.Errorf("TestGetLday.%s test fail, Error: %v", caseName, err);
        return
    }

    caseName = "Success"
    code := "sz000001"
    putMsg := newOneStockPutLdayRequest(code, 20161001, 300);

    t.Logf("KDataList Size: %d", len(putMsg.Data))
    err = putLday(putMsg)
    if err != nil {
        t.Errorf("TestGetLday.%s, putLday fail, Error: %v", caseName, err);
        return;
    }

    r6 := httptest.NewRequest("GET", "/api/v1/lday?type=csv&head=false&code=sz000001&from=20161001&to=20171001", nil)
    r6.Header.Add("Accept", "text/csv")
    w6 := httptest.NewRecorder()
    httptest.Do(w6, r6, apiV1GetLday)
//    t.Logf("%s ==> \nResponse: %d,\nBody:\n%s", caseName, w6.Code, w6.Body.String());
    if w6.Code != iris.StatusOK {
        t.Errorf("TestGetLday.%s test fail, statusCode: %d", caseName, w6.Code);
        return;
    }
    csvReader := csv.NewReader(w6.Body)
    err = checkCsvWithPutMsg(code, csvReader, putMsg.Data)
    if err != nil {
        t.Errorf("TestGetLday.%s, check reply csv fail: %s", caseName, err);
        return;
    }
}

func Test_PutLday(t *testing.T) {

    caseName := "BadProtoc"
    reqBin1 := make([]byte, 10)
    w1 := httptest.NewRecorder()
    r1 := httptest.NewRequest("PUT", "/api/v1/lday", bytes.NewReader(reqBin1))
    httptest.Do(w1, r1, apiV1PutLday)
//    t.Logf("%s ==> \nResponse: %d,\nBody:\n%s", caseName, w1.Code, w1.Body.String());
    err := checkResponseError(w1, StatusBadRequest, "cannot parse invalid wire-format data")
    if err != nil {
        t.Errorf("TestPutLday.%s, test fail, Error: %v", caseName, err);
    }

    caseName = "DataEmpty"
    reqMsg2 := &PutLdayRequest{
        Data: nil,
    }
    reqBin2, err := proto.Marshal(reqMsg2);
    if err != nil {
        t.Errorf("TestPutLday.%s, protobuf marshal LDayRequest fail, Error: %v", caseName, err);
        return;
    }
    w2 := httptest.NewRecorder()
    r2 := httptest.NewRequest("PUT", "/api/v1/lday", bytes.NewReader(reqBin2))
    httptest.Do(w2, r2, apiV1PutLday)
//    t.Logf("%s ==> \nResponse: %d,\nBody:\n%s", caseName, w2.Code, w2.Body.String());
    err = checkResponseError(w2, StatusBadRequest, "KData is empty")
    if err != nil {
        t.Errorf("TestPutLday.%s test fail, Error: %v", caseName, err);
        return
    }

    caseName = "MalforamtKData"
    var m3 ResponseModel
    l_kdata := make([]*KData, 100)
    reqMsg3 := &PutLdayRequest{
        Data: l_kdata,
    }
    reqBin3, err := proto.Marshal(reqMsg3);
    if err != nil {
        t.Errorf("TestPutLday.%s, protobuf marshal LdayRequest fail, Error: %v", caseName, err);
        return;
    }
    w3 := httptest.NewRecorder()
    r3 := httptest.NewRequest("PUT", "/api/v1/lday", bytes.NewReader(reqBin3))
    httptest.Do(w3, r3, apiV1PutLday)
    if (w3.Code != iris.StatusOK) {
        t.Errorf("TestPutLday.%s, Http Response StatusCode %d", caseName, w3.Code);
    }
    //t.Logf("%s ==> \nResponse: %d,\nBody:\n%s", caseName, w3.Code, w3.Body.String());
    err = json.Unmarshal(w3.Body.Bytes(), &m3);
    if err != nil {
        t.Errorf("TestPutLday.%s, Unmarshal response fail, Error: %s", caseName, err);
        return
    }
    if m3.Status != StatusOK {
        t.Errorf("TestPutLday.%s, Http Response Status %d", caseName, m3.Status);
        return
    }
    if m3.Message != "0 stock and 0 kdata updated" {
        t.Errorf("TestPutLday.%s, Check Update message fail: %s", caseName, m3.Message);
        return 
    }

    caseName = "Success"
    reqMsg4 := newOndDayPutLdayRequest(20160522, 300);
    reqBin4, err := proto.Marshal(reqMsg4);
    if err != nil {
        t.Errorf("TestPutLday.%s, protobuf marshal LdayRequest fail, Error: %v", caseName, err);
        return;
    }

    w4 := httptest.NewRecorder()
    r4 := httptest.NewRequest("PUT", "/api/v1/lday", bytes.NewReader(reqBin4))
    httptest.Do(w4, r4, apiV1PutLday)
//    t.Logf("%s ==> \nResponse: %d,\nBody:\n%s", caseName, w4.Code, w4.Body.String());

    ll := make([]*KData, 1)
    for _, k := range reqMsg4.Data {
        csvReader, err := getLday(k.Code, 20160522, 20160522)
        if err != nil {
            t.Errorf("TestPutLday.%s, get %s's Lday fail, Error: %s", k.Code, caseName, err)
            return 
        }
        ll[0] = k
        err = checkCsvWithPutMsg(k.Code, csvReader, ll)
        if err != nil {
            t.Errorf("TestPutLday.%s, check %s's reply csv fail: %s", k.Code, caseName, err);
            return;
        }
    }
}