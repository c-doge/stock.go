package gostk

import (
    "io"
	"os"
	"fmt"
	"encoding/csv"
    "github.com/c-doge/stock.go/base/utils"
)


//日期  开盘  最高  最低  收盘  成交量  成交额
func ReadLdayCsv(path string) ([]*KData, error) {
    list := make([]*KData, 0)
    fs, err := os.Open(path)
    if err != nil {
        return nil, err;
    }
    defer fs.Close()
    r := csv.NewReader(fs)
    total := 0
    for {
        record, err := r.Read();
        if err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        }

        t, err := utils.ParseTime("2006-01-02", record[0])
        if err != nil {
            return nil, fmt.Errorf("parse Date fail, error: %v", err);
        }
        k := &KData {
            Time:     t,
            Open:     float32(utils.ParseFloat(record[1])),
            Close:    float32(utils.ParseFloat(record[4])),
            High:     float32(utils.ParseFloat(record[2])),
            Low:      float32(utils.ParseFloat(record[3])),
            Volume:   utils.ParseFloat(record[5]),
            Turnover: utils.ParseFloat(record[6]),
        }
        list = append(list, k);
        total += 1;
    }
    return list, nil
}
