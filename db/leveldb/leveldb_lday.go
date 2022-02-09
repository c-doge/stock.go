package leveldb

import (
    "fmt"
    "time"
    "sort"
    "strings"
    "strconv"
    "github.com/golang/protobuf/proto"
//  "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/opt"
    "github.com/syndtr/goleveldb/leveldb/util"
    "github.com/syndtr/goleveldb/leveldb/errors"
    
    "github.com/c-doge/stock.go/gostk"
    "github.com/c-doge/stock.go/base/utils"
    "github.com/c-doge/stock.go/base/logger"
)

var dbLdayVersion string   = "0.0.1"

func checkLdayVersion() error {
    key := "version_lday"
    ver := ""
    b_ver, err := dbLday.Get([]byte(key), &opt.ReadOptions{})
    if err != nil {
        if err != errors.ErrNotFound {
            logger.Fatalf("[LevelDB] get lday version fail, err: %v", err);
            return err;
        }
        err := dbLday.Put([]byte(key), []byte(dbLdayVersion), nil)
        if err != nil {
            logger.Fatalf("[LevelDB] set lday version fail, err: %v", err);
            return err;
        }
        ver = dbLdayVersion;
    } else {
        ver = string(b_ver)
    }
    logger.Infof("[LevelDB] lday version: %s", ver);
    return nil
}

func mergeRecordLday(to *RecordLday, from *RecordLday) *RecordLday {
    if to.Code != from.Code {
        return nil;
    }
    list := make([]*KData, 0)
    var len_to int = len(to.Data)
    var len_from int = len(from.Data)
    var i int = 0
    var j int = 0
    for ;i < len_to && j < len_from; {
        if to.Data[i].Time == from.Data[j].Time { 
            list = append(list, to.Data[i])
            i += 1
            j += 1
        } else if to.Data[i].Time < from.Data[j].Time {
            list = append(list, to.Data[i])
            i += 1
        } else if to.Data[i].Time > from.Data[j].Time {
            list = append(list, from.Data[j])
            j += 1
        }
    }
    if i < len_to {
        list = append(list, to.Data[i:]...)
    }
    if j < len_from {
        list = append(list, from.Data[j:]...)
    }
    lday := new(RecordLday)
    lday.Data = list;
    return lday;
}

func convertStkKDataToRecordLday(code string, list []*gostk.KData) (*RecordLday, error) {
    if len(list) <= 0 {
        return nil, ErrorInvalidParameter;
    }
    var err error = nil;
    var ll = make([]*KData, 0)
    var year = int(list[0].Time.Year());
    for _, kdata := range list {
        if kdata.Time.Year() != year {
            err = ErrorInvalidParameter;
            break
        }
        data := &KData {
            Time:     utils.DateTimeToDecimalNum(kdata.Time),
            Open:     kdata.Open,
            Close:    kdata.Close,
            High:     kdata.High,
            Low:      kdata.Low,
            Volume:   kdata.Volume,
            Turnover: kdata.Turnover,
        }
        ll = append(ll, data)
    }
    if err != nil {
        return nil, err
    } else if len(ll) == 0 {
        return nil, ErrorNotResult
    }
    lday := &RecordLday{
                Code: code,
                Data: ll,
            }
    return lday, nil;
}

func convertRecordLdayToStkKDataList(lday *RecordLday, from, to time.Time) ([]*gostk.KData) {
    var list = make([]*gostk.KData, len(lday.Data));
    var i = 0;
    for _, v := range lday.Data {
        
        t := utils.DecimalNumToDateTime(v.GetTime());
        if  t.Before(from) {
            continue;
        } else if t.After(to) {
            break;
        }
        kdata := &gostk.KData {
            Time:       t,
            Open:       v.GetOpen(),
            Close:      v.GetClose(),
            High:       v.GetHigh(),
            Low:        v.GetLow(),
            Volume:     v.GetVolume(),
            Turnover:   v.GetTurnover(),
        }

        list[i] = kdata;
        i += 1
    }
    return list[:i];
}

func getRecordLday(code string, year int) (*RecordLday, error) {
    key := fmt.Sprintf("%s_%d", code, year); 
    bin, err := dbLday.Get([]byte(key), &opt.ReadOptions{})
    if err != nil {
        return nil, err;
    }
    lday := &RecordLday{};
    err = proto.Unmarshal(bin, lday)
    if err != nil {
        return nil, err;
    }
    logger.Debugf("[LevelDB] getlday key(%v), len(lday.Data): %d", key, len(lday.Data))
    return lday, nil;
}

func putRecordLday(code string, year int, lday *RecordLday) error {
    bin, err := proto.Marshal(lday);
    if err != nil {
        return err;
    }
    key := fmt.Sprintf("%s_%d", code, year); 
    err = dbLday.Put([]byte(key), bin, nil)
    if err != nil {
        return err;
    }
    logger.Debugf("[LevelDB] putlday key(%v), len(lday.Data): %d", key, len(lday.Data))
    return nil
}


func PutLday(code string, list []*gostk.KData) error {
    sort.SliceStable(list, func(i, j int) bool {
        return list[i].Time.Before(list[j].Time)
    })
    var err error = nil
    var year int = -1;
    var lday *RecordLday = nil;

    for _, v := range list {
        d := &KData {
            Time:      utils.DateTimeToDecimalNum(v.Time),
            Open:      v.Open,
            Close:     v.Close,
            High:      v.High,
            Low:       v.Low,
            Volume:    v.Volume,
            Turnover:  v.Turnover,
        }
        if year != -1 && year != v.Time.Year() {
            lday2, err := getRecordLday(code, year)
            if err == nil && lday2 != nil {
                mergeRecordLday(lday, lday2)
            } 
            err = putRecordLday(code, year, lday);
            if err != nil {
                return err
            }
            lday = nil;
        } 
        if lday == nil {
            lday = new(RecordLday)
            lday.Data = make([]*KData, 0)
        }
        lday.Data = append(lday.Data, d);
        year = v.Time.Year()
    }
    if lday != nil {
        err = putRecordLday(code, year, lday);
        if err != nil {
            return err
        }
        lday = nil;
    }
    return err;
}

func GetLday(code string, from, to time.Time) ([]*gostk.KData, error) {
    if from.Year() > to.Year() {
        return nil, ErrorInvalidParameter;
    }
    if !gostk.IsStockCodeValid(code) {
        return nil, ErrorInvalidParameter;
    }
    var err error = nil;
    var start_ts time.Time
    var end_ts time.Time
    var result []*gostk.KData = make([]*gostk.KData, 0)

    startKey := fmt.Sprintf("%s_%d", code, from.Year());
    limitKey := fmt.Sprintf("%s_%d", code, to.Year() + 1);
    iter := dbLday.NewIterator(&util.Range{Start: []byte(startKey), Limit: []byte(limitKey)}, nil)

    for iter.Next() {
        key := iter.Key();
        i := strings.LastIndex(string(key), "_")
        year, _ := strconv.Atoi(string(key[i + 1:]))
        if !iter.Valid() {
            return nil, fmt.Errorf("Iterator Invalied");
        } else if iter.Error() != nil {
            return nil, iter.Error();
        }
        v := iter.Value();
        lday := &RecordLday{};
        err = proto.Unmarshal(v, lday)
        if err != nil {
            return nil, err;
        }
        start_ts = utils.EarlyDate
        end_ts = utils.FutureDate
        if from.Year() == year {
            start_ts = from
        }
        if to.Year() == year {
            end_ts = to
        }
        list := convertRecordLdayToStkKDataList(lday, start_ts, end_ts);
        if len(list) > 0 {
            result = append(result, list...)
        }
    }
    iter.Release()
    err = iter.Error()

    if err == nil && len(result) == 0 {
        result = nil;
        err = ErrorNotResult;
    }
    return result, err;
}