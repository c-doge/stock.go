package api

import (
    "io"
    "fmt"
    "math"
    "bytes"
    "strings"
    "testing"
    "math/rand"
    "encoding/csv"
    "encoding/json"
    htest "net/http/httptest"

    "github.com/kataras/iris/v12"
    "github.com/kataras/iris/v12/httptest"
    "github.com/golang/protobuf/proto"
    "github.com/c-doge/stock.go/base/utils"
    "github.com/c-doge/stock.go/base/logger"

)

func equalVData(d1, d2 *VData) bool {
    if d1.Date != d2.Date {
        logger.Debugf("1: %d, %d", d1.Date , d2.Date);
        return false
    } else if math.Abs(d1.PreTradable - d2.PreTradable) > 0.00001 {
        return false
    } else if math.Abs(d1.PostTradable - d2.PostTradable) > 0.00001 {
        return false
    } else if math.Abs(d1.PreTotal - d2.PreTotal) > 0.00001 {
        return false
    } else if math.Abs(d1.PostTotal - d2.PostTotal) > 0.00001 {
        return false
    }
    return true;
}
func equalXData(d1, d2 *XData) bool {
    if d1.Date != d2.Date {
        return false
    } else if abs32(d1.AllotVolume, d2.AllotVolume) > 0.00001 {
        return false
    } else if abs32(d1.AllotPrice, d2.AllotPrice) > 0.00001 {
        return false
    } else if abs32(d1.BonusVolume, d2.BonusVolume) > 0.00001 {
        return false
    } else if abs32(d1.BonusPrice, d2.BonusPrice) > 0.00001 {
        return false
    }
    return true;
}

func checkVolResponseError(w *htest.ResponseRecorder, status int, errMsg string) error {

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

func checkVolCsvResponse(code string, r *csv.Reader, ll []*VData) error {
    index := 0;
    for {
        record, err := r.Read();
        if err == io.EOF {
            break
        } else if err != nil {
            return fmt.Errorf("Read CSV Fail, Error: %s", err);
        }
        if index >= len(ll) {
            return fmt.Errorf("the VData got is more than put, index: %d", index);
        }
        t, err := utils.ParseTime("2006-01-02", record[0])
        if err != nil {
            return fmt.Errorf("parse Date fail, error: %v", err);
        }
        v := &VData {
            Date:            utils.DateTimeToDecimalNum(t),
            PreTradable:     parseFloat(record[1]),
            PostTradable:    parseFloat(record[2]),
            PreTotal:        parseFloat(record[3]),
            PostTotal:       parseFloat(record[4]),
        }
        if !equalVData(v, ll[index]) {
            return fmt.Errorf("csv[%d] != putReq.Data[%d], Index:%v", v.Date, ll[index].Date, index);
        }
        index += 1;
    }
    if index != len(ll) {
        return fmt.Errorf("lines of csv(%d) not equal to the size of putReq.Data", index);
    }
    return nil
}

func checkXdrCsvResponse(code string, r *csv.Reader, ll []*XData) error {
    index := 0;
    for {
        record, err := r.Read();
        if err == io.EOF {
            break
        } else if err != nil {
            return fmt.Errorf("Read CSV Fail, Error: %s", err);
        }
        if index >= len(ll) {
            return fmt.Errorf("the XData got is more than put, index: %d", index);
        }
        t, err := utils.ParseTime("2006-01-02", record[0])
        if err != nil {
            return fmt.Errorf("parse Date fail, error: %v", err);
        }
        x := &XData {
            Date:             utils.DateTimeToDecimalNum(t),
            AllotPrice:       float32(parseFloat(record[1])),
            AllotVolume:      float32(parseFloat(record[2])),
            BonusPrice:       float32(parseFloat(record[3])),
            BonusVolume:      float32(parseFloat(record[4])),
        }
        if !equalXData(x, ll[index]) {
            return fmt.Errorf("csv[%d] != putReq.Data[%d], Index:%v", x.Date, ll[index].Date, index);
        }
        index += 1;
    }
    if index != len(ll) {
        return fmt.Errorf("lines of csv(%d) not equal to the size of putReq.Data", index);
    }
    return nil
}

func newPutVolRequest(code string, dateNum uint32, n int) *PutVolRequest {
    dateTime := utils.DecimalNumToDateTime(dateNum)
    l := make([]*VData, n)
    for i:= 0; i < n; i++ {
        pre := rand.Float64() * 1000000
        post := pre + rand.Float64() * 1000
        l[i] = &VData {
            Date:          utils.DateTimeToDecimalNum(dateTime),
            PreTradable:   pre,
            PreTotal:      pre + rand.Float64() * 1000,
            PostTradable:  post,
            PostTotal:     post + rand.Float64() * 1000,    
        }
        logger.Debugf("i:%d, %v", i, l[i])
        dateTime = dateTime.AddDate(0, 0 ,1)
    }
    req := &PutVolRequest {
        Code: code,
        Data: l,
    }
    return req
}

func newPutXDRRequest(code string, dateNum uint32, n int) *PutXDRRequest {
    dateTime := utils.DecimalNumToDateTime(dateNum)
    l := make([]*XData, n)
    for i:= 0; i < n; i++ {
        l[i] = &XData {
            Date:          utils.DateTimeToDecimalNum(dateTime),
            AllotVolume:   rand.Float32() * 100,
            AllotPrice:    rand.Float32() * 10,
            BonusVolume:   rand.Float32() * 100,
            BonusPrice:    rand.Float32() * 10,   
        }
        dateTime = dateTime.AddDate(0, 0 ,1)
    }
    req := &PutXDRRequest {
        Code: code,
        Data: l,
    }
    return req
}

func Test_volGetForBadType(t *testing.T) {
    // bad type
    r := httptest.NewRequest("GET", "/api/v1/vol?type=csx&code=sh000001&from=20120101&to=20120131", nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiV1GetVolumeList)
   // t.Logf("Test volGetForBadType ==> \nResponse: %d,\nBody:%s", w.Code, w.Body.String());
    err := checkVolResponseError(w, StatusBadRequest, "Accept Type Unknown")
    if err != nil {
        t.Errorf("Test volGetForBadType fail, Error: %v", err);
        return
    }
}

func Test_volGetForBadCode(t *testing.T) {
    // bad code
    r := httptest.NewRequest("GET", "/api/v1/vol?type=csv&code=sa000001&from=20120101&to=20120131", nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiV1GetVolumeList)
    //t.Logf("Test volGetForBadCode ==> \nResponse: %d,\nBody:\n%s", w.Code, w.Body.String());
    err := checkVolResponseError(w, StatusBadRequest, "Stock Code Unknown")
    if err != nil {
        t.Errorf("Test volGetForBadCode fail, Error: %v", err);
        return
    }
}

func Test_volGetForBadFromTime(t *testing.T) {
    r := httptest.NewRequest("GET", "/api/v1/vol?type=csv&code=sh000001&from=120120101&to=20120131", nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiV1GetVolumeList)
//    t.Logf("Test volGetForBadFromTime ==> \nResponse: %d,\nBody:\n%s", w.Code, w.Body.String());
    err := checkVolResponseError(w, StatusBadRequest, "Parse 120120101 fail")
    if err != nil {
        t.Errorf("Test volGetForBadFromTime fail, Error: %v", err);
        return
    }
}
func Test_volGetForBadToTime(t *testing.T) {
    r := httptest.NewRequest("GET", "/api/v1/vol?type=csv&code=sh000001&from=20120101&to=210120131", nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiV1GetVolumeList)
//    t.Logf("Test volGetForBadToTime ==> \nResponse: %d,\nBody:\n%s", w.Code, w.Body.String());
    err := checkVolResponseError(w, StatusBadRequest, "Parse 210120131 fail")
    if err != nil {
        t.Errorf("Test volGetForBadToTime fail, Error: %v", err);
        return
    }
}
func Test_volGetForFromTimeLaterThanTo(t *testing.T) {
    r := httptest.NewRequest("GET", "/api/v1/vol?type=csv&code=sh000001&from=20120101&to=20110131", nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiV1GetVolumeList)
//    t.Logf("Test volGetForFromTimeLaterThanTo ==> \nResponse: %d,\nBody:\n%s", w.Code, w.Body.String());
    err := checkVolResponseError(w, StatusBadRequest, "From time is later than To time")
    if err != nil {
        t.Errorf("Test volGetForFromTimeLaterThanTo fail, Error: %v", err);
        return
    }   
}

func Test_volGetAndPut(t *testing.T) {
    reqMsg := newPutVolRequest("sz000010", 20161022, 10)
    reqBin, err := proto.Marshal(reqMsg);
    if err != nil {
        t.Errorf("Test volGetAndPut, protobuf marshal LdayRequest fail, Error: %v", err);
        return;
    }
    w1 := httptest.NewRecorder()
    r1 := httptest.NewRequest("PUT", "/api/v1/vol", bytes.NewReader(reqBin))
    httptest.Do(w1, r1, apiV1PutVolumeList)
//    t.Logf("Test volGetAndPut ==> \nPut Response: %d,\nBody:\n%s", w1.Code, w1.Body.String());
    err = checkVolResponseError(w1, StatusOK, "")
    if err != nil {
        t.Errorf("Test volGetAndPut fail, %s", err);
        return;
    }
    w2 := httptest.NewRecorder()
    r2 := httptest.NewRequest("GET", "/api/v1/vol?type=csv&code=sz000010&head=true", nil)
    r2.Header.Add("Accept", "text/csv")
    httptest.Do(w2, r2, apiV1GetVolumeList)
//    t.Logf("Test volGetAndPut ==> \nGet Response: %d,\nBody:\n%s", w2.Code, w2.Body.String());
    if w2.Code != iris.StatusOK {
        t.Errorf("Test volGetAndPut, Response Status Code: %d", w2.Code)
    }
    reader := csv.NewReader(w2.Body)
    _, err = reader.Read()
    if err != nil {
        t.Errorf("Test xdrGetAndPut, read csv head fail, Error: %v", err)
        return
    }
    err = checkVolCsvResponse("sz000010", reader, reqMsg.Data)
    if err != nil {
        t.Errorf("Test volGetAndPut, check csv response fail, %s", err);
    }
}

///////////////////////////


func Test_xdrGetForBadType(t *testing.T) {
    // bad type
    r := httptest.NewRequest("GET", "/api/v1/xdr?type=csx&code=sh000001&from=20120101&to=20120131", nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiV1GetXDRList)
    //t.Logf("Test xdrGetForBadType ==> \nResponse: %d,\nBody:%s", w.Code, w.Body.String());
    err := checkVolResponseError(w, StatusBadRequest, "Accept Type Unknown")
    if err != nil {
        t.Errorf("Test xdrGetForBadType fail, Error: %v", err);
        return
    }
}

func Test_xdrGetForBadCode(t *testing.T) {
    // bad code
    r := httptest.NewRequest("GET", "/api/v1/xdr?type=csv&code=sa000001&from=20120101&to=20120131", nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiV1GetXDRList)
//    t.Logf("Test xdrGetForBadCode ==> \nResponse: %d,\nBody:\n%s", w.Code, w.Body.String());
    err := checkVolResponseError(w, StatusBadRequest, "Stock Code Unknown")
    if err != nil {
        t.Errorf("Test xdrGetForBadCode fail, Error: %v", err);
        return
    }
}

func Test_xdrGetForBadFromTime(t *testing.T) {
    r := httptest.NewRequest("GET", "/api/v1/xdr?type=csv&code=sh000001&from=120120101&to=20120131", nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiV1GetXDRList)
//    t.Logf("Test xdrGetForBadFromTime ==> \nResponse: %d,\nBody:\n%s", w.Code, w.Body.String());
    err := checkVolResponseError(w, StatusBadRequest, "Parse 120120101 fail")
    if err != nil {
        t.Errorf("Test xdrGetForBadFromTime fail, Error: %v", err);
        return
    }
}
func Test_xdrGetForBadToTime(t *testing.T) {
    r := httptest.NewRequest("GET", "/api/v1/xdr?type=csv&code=sh000001&from=20120101&to=210120131", nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiV1GetXDRList)
//    t.Logf("Test xdrGetForBadToTime ==> \nResponse: %d,\nBody:\n%s", w.Code, w.Body.String());
    err := checkVolResponseError(w, StatusBadRequest, "Parse 210120131 fail")
    if err != nil {
        t.Errorf("Test xdrGetForBadToTime fail, Error: %v", err);
        return
    }
}

func Test_xdrGetForFromTimeLaterThanTo(t *testing.T) {
    r := httptest.NewRequest("GET", "/api/v1/xdr?type=csv&code=sh000001&from=20120101&to=20110131", nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiV1GetXDRList)
//    t.Logf("Test xdrGetForFromTimeLaterThanTo ==> \nResponse: %d,\nBody:\n%s", w.Code, w.Body.String());
    err := checkVolResponseError(w, StatusBadRequest, "From time is later than To time")
    if err != nil {
        t.Errorf("Test xdrGetForFromTimeLaterThanTo fail, Error: %v", err);
        return
    }   
}

func Test_xdrGetAndPut(t *testing.T) {
    reqMsg := newPutXDRRequest("sz000010", 20161022, 10)
    reqBin, err := proto.Marshal(reqMsg);
    if err != nil {
        t.Errorf("Test xdrGetAndPut, protobuf marshal PutXDRRequest fail, Error: %v", err);
        return;
    }
    w1 := httptest.NewRecorder()
    r1 := httptest.NewRequest("PUT", "/api/v1/xdr", bytes.NewReader(reqBin))
    httptest.Do(w1, r1, apiV1PutXDRList)
//    t.Logf("Test volGetAndPut ==> \nPut Response: %d,\nBody:\n%s", w1.Code, w1.Body.String());
    err = checkVolResponseError(w1, StatusOK, "")
    if err != nil {
        t.Errorf("Test xdrGetAndPut fail, %s", err);
        return;
    }
    w2 := httptest.NewRecorder()
    r2 := httptest.NewRequest("GET", "/api/v1/xdr?type=csv&code=sz000010&head=true", nil)
    r2.Header.Add("Accept", "text/csv")
    httptest.Do(w2, r2, apiV1GetXDRList)
//    t.Logf("Test xdrGetAndPut ==> \nGet Response: %d,\nBody:\n%s", w2.Code, w2.Body.String());
    if w2.Code != iris.StatusOK {
        t.Errorf("Test xdrGetAndPut, Response Status Code: %d", w2.Code)
    }
    reader := csv.NewReader(w2.Body)
    _, err = reader.Read()
    if err != nil {
        t.Errorf("Test xdrGetAndPut, read csv head fail, Error: %v", err)
        return
    }
    err = checkXdrCsvResponse("sz000010", reader, reqMsg.Data)
    if err != nil {
        t.Errorf("Test xdrGetAndPut, check csv response fail, %s", err);
    }
}
