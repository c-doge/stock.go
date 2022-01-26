package utils

import (
    "os"
    "fmt"
    "testing"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "path/filepath"
    "github.com/c-doge/stock.go/base/logger"
)

var testPath string = "../../test"

func checkUnzipFiles(dstPath string) error {
    names := []string{
        "a",
        "b",
        "c",
        "d",
    }

    for _, name := range names {
        p := filepath.Join(dstPath, name)
        _, err := os.Stat(p)
        if err != nil {
            return err;
        }
    }
    return nil
}

func Test_unzip(t *testing.T) {
    zipFile := filepath.Join(testPath, "z.zip")
    dstPath := filepath.Join(testPath, "z")
    f, err := os.Open(zipFile)
    defer f.Close()
    if err != nil {
        t.Errorf("open lday v2 fail,  Error: %v", err)
        return
    }
    fi, err := f.Stat()
    if err != nil {
        t.Errorf("zip file stats fail, Error: %v", err)
        return;
    }
    
    err = Mkdir(dstPath)
    if err != nil {
        t.Errorf("mkdir %s fail, Error: %v", dstPath, err)
    }
    err = Unzip(f, fi.Size(), dstPath)
    if err != nil {
        t.Errorf("Unzip fail,  Error: %v", err)
        return
    }
    err = checkUnzipFiles(dstPath)
    if err != nil {
        t.Errorf("checkUnzipFiles fail, Error: %v", err)
    }
    os.RemoveAll(dstPath)
}

func Test_DownloadOnly(t *testing.T) {
    zipFile := filepath.Join(testPath, "z.zip")
    dstPath := filepath.Join(testPath, "z")

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var content []byte    
        s, err := ioutil.ReadFile(zipFile)
        if err != nil {
            content = []byte("404")
        } else {
            content = []byte(s)
            w.Header().Set("Content-Type", "application/zip")
            
        }
        fmt.Fprintf(w, "%s", content)
    }))
    defer ts.Close()
    url := ts.URL
    t.Logf("url: %s", url);

    err := Mkdir(dstPath)
    if err != nil {
        t.Errorf("mkdir %s fail, Error: %v", dstPath, err)
    }
    dstFile := filepath.Join(dstPath, "d.zip")
    err = Download(url, dstFile, false)
    if err != nil {
        t.Errorf("Download fail, Error: %v", err);
    }

    fi, err := os.Stat(dstFile);
    if err != nil {
        t.Errorf("check ZipFile fail, Error: %v", err);
    } else if fi.Size() != 586 {
        t.Errorf("check ZipFile fail, file size: %d != %d", fi.Size(), 586);
    }
    os.RemoveAll(dstPath)
}

func Test_DownloadUnzip(t *testing.T) {
    zipFile := filepath.Join(testPath, "z.zip")
    dstPath := filepath.Join(testPath, "z")

    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var content []byte    
        s, err := ioutil.ReadFile(zipFile)
        if err != nil {
            content = []byte("404")
        } else {
            content = []byte(s)
            w.Header().Set("Content-Type", "application/zip")
            
        }
        fmt.Fprintf(w, "%s", content)
    }))
    defer ts.Close()
    url := ts.URL
    t.Logf("url: %s", url);

    err := Mkdir(dstPath)
    if err != nil {
        t.Errorf("mkdir %s fail, Error: %v", dstPath, err)
    }

    err = Download(url, dstPath, true)
    if err != nil {
        t.Errorf("Download fail, Error: %v", err);
    }

    err = checkUnzipFiles(dstPath)
    if err != nil {
        t.Errorf("checkUnzipFiles fail, Error: %v", err)
    }
    os.RemoveAll(dstPath)
}

func Test_TimeData(t *testing.T) {
    date := uint32(20120322)
    t1 := DateToTime(date)
    d1 := TimeToDate(t1);
    if d1 != date {
        t.Errorf("%d not equal to %d", d1, 20120322)
    }
}

func TestMain(m *testing.M) {
    
    logger.New("Debug", "", "stock.go/base/utils")

    logger.Info("stock.go base/utils test start >>>")
    m.Run()
    logger.Info("stock.go base/utils test stop >>>")

   // 

}