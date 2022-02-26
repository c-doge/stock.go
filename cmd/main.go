package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/c-doge/stock.go/db"
    "github.com/c-doge/stock.go/api"
    "github.com/c-doge/stock.go/base"
    "github.com/c-doge/stock.go/base/logger"
    
)

var help bool
var configFile string

func usage() {
    fmt.Printf("stock.go version: %s\r\n", version)
    fmt.Printf("Usage: stock.go [-ch]\r\n")
    fmt.Printf("           -h print this message\r\n")
    fmt.Printf("           -c config file path\r\n")
}

func init() {
    flag.BoolVar(&help,             "h", false,                "show this help")
    flag.StringVar(&configFile,     "c", "./gostock.yaml",       "config file path")
    flag.Usage = usage
}


func main() {

    flag.Parse()
    if help {
        usage()
        os.Exit(0)
    }
    err := base.Start(configFile);
    if err != nil {
        panic(err);
    }

    logger.Infof("version:     %s\n", _version);
    logger.Infof("git branch:  %s\n", _gitBranch);
    logger.Infof("build time:  %s\n", _buildTime);

    db.Start();
    api.Init(true)


    db.Stop();
    base.Stop();

}