package api

import (
//    "os"
    "fmt"
    "time"
    "bytes"
    "testing"
    "encoding/json"

    "github.com/golang/protobuf/proto"
    "github.com/kataras/iris/v12/httptest"

    "github.com/c-doge/stock.go/gostk"
    "github.com/c-doge/stock.go/base/logger"
    "github.com/c-doge/stock.go/base/setting"
)

var configFile string = "../cmd/gostock-ut.yaml"

/*
func newPutLdayRequest() *PutLdayRequest {
    k1 := &KData {
        Code:    "sh000001",
        Time:     20120115,
        Open:     10.1,
        Close:    10.2,
        High:     15.1,
        Low:      9.1,
        Volume:   1000.1,
        Turnover: 12345.6,

    }
    k2 := &KData {
        Code:    "sz000002",
        Time:     20120115,
        Open:     5.1,
        Close:    6.2,
        High:     7.1,
        Low:      4.1,
        Volume:   1000.1,
        Turnover: 54321.0,
    }
    l := make([]*KData, 0)
    l = append(l, k1, k2);
    req := &PutLdayRequest {
        Data: l,
    }
    return req
}
*/
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

func newPutLdayRequest(n int) *PutLdayRequest {
    max := 100;
    dateNum := 20120115
    dateTime := utils.DecimalNumToDateTime(dateNum)
    l := make([]*KData, max)
    for i:= 0: i < max; i++ {
        code = fmt.Sprintf("sh%6d", 600000+i+1)
        k := newStkKData(dateTime)
        kdata := &KData {
            Code:     code,
            Time:     dateNum,
            Open:     k.Open,
            Close:    k.Close,
            High:     k.High,
            Low:      k.Low,
            Volume:   k.Volume,
            Turnover: k.Turnover,
        }
        l = append(l, kdata)
    }
    req := &PutLdayRequest {
        Data: l,
    }
    return req
}
func putLday() {
    reqMsg := newPutLdayRequest();
    reqBin, err := proto.Marshal(reqMsg);
    if err != nil {
        t.Errorf("protobuf marshal LDayRequest fail, Error: %v", err);
        return;
    }
    w := httptest.NewRecorder()
    r := httptest.NewRequest("PUT", "/api/lday", bytes.NewReader(reqBin))
    httptest.Do(w, r, apiPutLday)
}

func Test_GetLday(t *testing.T) {

    r := httptest.NewRequest("GET", "/api/lday?type=csv&code=sh000001&from=20120101&tg=20120131", nil)
    r.Header.Add("Accept", "text/csv")
    w := httptest.NewRecorder()
    httptest.Do(w, r, apiGetLday)
     t.Logf("Response: %d \n, Body: %s", w.Code, w.Body.String());
}

func Test_PutLday(t *testing.T) {
    reqMsg := newPutLdayRequest();
    reqBin, err := proto.Marshal(reqMsg);
    if err != nil {
        t.Errorf("protobuf marshal LDayRequest fail, Error: %v", err);
        return;
    }
    w1 := httptest.NewRecorder()
    r1 := httptest.NewRequest("PUT", "/api/lday", bytes.NewReader(reqBin))
    httptest.Do(w1, r1, apiPutLday)

    t.Logf("Response: %d \n, Body: %s", w1.Code, w1.Body.String());
}

func Test_GetVolume(t *testing.T) {

}

func Test_PutVolume(t *testing.T) {

}

func Test_CheckHealth(t *testing.T) {

    type Message struct {
        Timestamp int64 `json:"timestamp"`
        Status    int32 `json:"status"`
    }
    var m1 Message;
    w1 := httptest.NewRecorder()
    r1 := httptest.NewRequest("GET", "/api/health", nil)
    t1 := time.Now().Unix()
    httptest.Do(w1, r1, apiHealth)
    err := json.Unmarshal(w1.Body.Bytes(), &m1);
    if err != nil {
        t.Errorf("TestHealth, Unmarshal response fail, Error: %s", err);
    }
    if m1.Status != httptest.StatusOK {
         t.Errorf("TestHealth status Code Not Equal to %d", httptest.StatusOK);
    }
    if m1.Timestamp != t1 {
        t.Errorf("TestHealth timestamp Code Not Equal to %d", t1);
    }

    time.Sleep(time.Duration(2) * time.Second);

    var m2 Message;
    r2 := httptest.NewRequest("POST", "/api/health", nil)
    w2 := httptest.NewRecorder()
    t2 := time.Now().Unix()
    httptest.Do(w2, r2, apiHealth)
    err = json.Unmarshal(w2.Body.Bytes(), &m2);
    if err != nil {
        t.Errorf("TestHealth, Unmarshal response fail, Error: %s", err);
    }
    if m2.Status != httptest.StatusOK {
         t.Errorf("TestHealth status Code Not Equal to %d", httptest.StatusOK);
    }
    if m2.Timestamp != t2 {
        t.Errorf("TestHealth timestamp Code Not Equal to %d", t2);
    }
}


func TestMain(m *testing.M) {

    s := setting.Get()
    s.Web.Addr = "127.0.0.1:9000"
    s.LevelDB.Path = "../../test/leveldb"

    fmt.Printf("setting:\n\r%v\n", s)

    logger.New("Debug", "", "stock.go/api")


    logger.Info("stock.go api test start >>>")
    Init(false);
    db.Start();
    m.Run()
    db.Stop();

    logger.Info("stock.go api test stop >>>")

   // 

}