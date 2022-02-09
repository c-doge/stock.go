package api

import (
    "os"
    "fmt"
    "time"
    "strings"
    "testing"
    "encoding/json"

    "github.com/kataras/iris/v12/httptest"

    "github.com/c-doge/stock.go/db"
    "github.com/c-doge/stock.go/base/utils"
    "github.com/c-doge/stock.go/base/logger"
    "github.com/c-doge/stock.go/base/setting"
)

var configFile string = "../cmd/gostock-ut.yaml"

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
    r1 := httptest.NewRequest("GET", "/api/v1/health", nil)
    t1 := time.Now().Unix()
    httptest.Do(w1, r1, apiV1Health)
    err := json.Unmarshal(w1.Body.Bytes(), &m1);
    if err != nil {
        t.Errorf("TestHealth, Unmarshal response fail, Error: %s", err);
        return
    }
    if m1.Status != httptest.StatusOK {
         t.Errorf("TestHealth status Code Not Equal to %d", httptest.StatusOK);
         return
    }
    if m1.Timestamp != t1 {
        t.Errorf("TestHealth timestamp Code Not Equal to %d", t1);
        return
    }

    time.Sleep(time.Duration(2) * time.Second);

    var m2 Message;
    r2 := httptest.NewRequest("POST", "/api/v1/health", nil)
    w2 := httptest.NewRecorder()
    t2 := time.Now().Unix()
    httptest.Do(w2, r2, apiV1Health)
    err = json.Unmarshal(w2.Body.Bytes(), &m2);
    if err != nil {
        t.Errorf("TestHealth, Unmarshal response fail, Error: %s", err);
        return;
    }
    if m2.Status != httptest.StatusOK {
        t.Errorf("TestHealth status Code Not Equal to %d", httptest.StatusOK);
        return;
    }
    if m2.Timestamp != t2 {
        t.Errorf("TestHealth timestamp Code Not Equal to %d", t2);
        return;
    }
}

func TestMain(m *testing.M) {

    s := setting.Get()
    s.Web.Addr = "127.0.0.1:9000"
    s.LevelDB.Path = "../test/leveldb"

    fmt.Printf("setting:\n\r%v\n", s)
    err := utils.Mkdir(s.LevelDB.Path);
    if err != nil {
        panic(err)
    }
    logger.New("Debug", "", "stock.go/api")
    logger.Info("stock.go api test start >>>")
    Init(false);
    db.Start();
    m.Run()
    db.Stop();

    if strings.Contains(s.LevelDB.Path, "/test/leveldb") {
        logger.Infof("Remove LevelDB test folder")
        os.RemoveAll(s.LevelDB.Path)
    }
    logger.Info("stock.go api test stop >>>")

   // 

}