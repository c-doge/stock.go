package db
import (

    "time"
    "github.com/c-doge/stock.go/gostk"
    "github.com/c-doge/stock.go/base/utils"
    "github.com/c-doge/stock.go/base/logger"    
    "github.com/c-doge/stock.go/base/setting"
)

type DBHelper interface {
    Start(dbPath string) error 
    Stop()
    // Lday
    PutLday(code string, list []*gostk.KData) error 
    GetLday(code string, from, to time.Time) ([]*gostk.KData, error)
    //Volume Change
    PutVolumeList(code string, list []*gostk.VData) error
    GetVolumeList(code string, from, to time.Time) ([]*gostk.VData, error)
    //XDR
    PutXDRList(code string, list []*gostk.XData) error
    GetXDRList(code string, from, to time.Time) ([]*gostk.XData, error)
}

var dbHelper DBHelper = nil

func Start() {
    if dbHelper != nil {
        logger.Warnf("DB helper has started!");
        return;
    }
    dbHelper = &leveldbHelper{};
    if dbHelper == nil {
        panic("dbHelper have not been set!");
    }
    s := setting.Get();
    if s == nil  {
        panic("db path have not been set!")
    }
    // for LevelDB
    leveldbPath := s.LevelDB.Path;
    if !utils.PathExists(leveldbPath) {
        panic("leveldb file path is not exist!")
    }
    err := dbHelper.Start(leveldbPath);
    if err != nil {
        panic(err);
    }
}

func Stop() {
    if dbHelper == nil {
        panic("dbHelper have not been set!");
    }
    dbHelper.Stop()
}

func PutLday(code string, list []*gostk.KData) error {
    if dbHelper == nil {
        panic("dbHelper have not been set!");
    }
    return dbHelper.PutLday(code, list)
}

func GetLday(code string, from, to time.Time) ([]*gostk.KData, error) {
    if dbHelper == nil {
        panic("dbHelper have not been set!");
    }
    return dbHelper.GetLday(code, from, to)
}

//Volume Change
func PutVolumeList(code string, list []*gostk.VData) error {
    if dbHelper == nil {
        panic("dbHelper have not been set!");
    }
    return dbHelper.PutVolumeList(code, list)
}
func GetVolumeList(code string, from, to time.Time) ([]*gostk.VData, error) {
    if dbHelper == nil {
        panic("dbHelper have not been set!");
    }
    return dbHelper.GetVolumeList(code, from, to)
}

//XDR
func PutXDRList(code string, list []*gostk.XData) error {
    if dbHelper == nil {
        panic("dbHelper have not been set!");
    }
    return dbHelper.PutXDRList(code, list)
}
func GetXDRList(code string, from, to time.Time) ([]*gostk.XData, error) {
    if dbHelper == nil {
        panic("dbHelper have not been set!");
    }
    return dbHelper.GetXDRList(code, from, to)
}