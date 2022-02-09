package leveldb

import (
    "os"
    "testing"
    "strings"
    
    "github.com/c-doge/stock.go/base/utils"
    "github.com/c-doge/stock.go/base/logger"
)

func abs32(f1, f2 float32) float32 {
    var v float32 = f1 - f2;
    if v > 0 {
        return v;
    }
    return 0-v;
}

func TestMain(m *testing.M) {

    leveldbPath := "../../test/leveldb"

    logger.New("Debug", "", "stock.go/db/leveldb")

    logger.Info("stock.go leveldb test start >>>")
    err := utils.Mkdir(leveldbPath);
    if err != nil {
        panic(err)
    }
    Start(leveldbPath);
    m.Run()
    Stop();
    if strings.Contains(leveldbPath, "/test/leveldb") {
        logger.Infof("Remove LevelDB test folder")
        os.RemoveAll(leveldbPath)
    }
    logger.Info("stock.go leveldb test stop >>>")

}