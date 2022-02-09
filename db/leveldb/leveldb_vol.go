package leveldb

import (
    "fmt"
    "time"
    "sort"
    
    "github.com/golang/protobuf/proto"
//  "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/opt"
    "github.com/syndtr/goleveldb/leveldb/errors"

    "github.com/c-doge/stock.go/gostk"
    "github.com/c-doge/stock.go/base/utils"
    "github.com/c-doge/stock.go/base/logger"
)

var dbVolVersion string   = "0.0.1"

func checkVolVersion() error {
    key := "version_vol"
    ver := ""
    b_ver, err := dbBase.Get([]byte(key), &opt.ReadOptions{})
    if err != nil {
        if err != errors.ErrNotFound {
            logger.Fatalf("[LevelDB] get vol version fail, err: %v", err);
            return err;
        }
        err := dbBase.Put([]byte(key), []byte(dbLdayVersion), nil)
        if err != nil {
            logger.Fatalf("[LevelDB] set vol version fail, err: %v", err);
            return err;
        }
        ver = dbLdayVersion;
    } else {
        ver = string(b_ver)
    }
    logger.Infof("[LevelDB] vol version: %s", ver);
    return nil
}

func convertStkVDataToRecordVol(code string, list []*gostk.VData) (*RecordVol, error) {
    if len(list) <= 0 {
        return nil, ErrorInvalidParameter;
    }
    var ll = make([]*VData, len(list))
    for i, v := range list {
        vv := &VData {
            Date:     utils.DateTimeToDecimalNum(v.Date),
            PreTradable:  v.PreTradable,
            PreTotal:     v.PreTotal,
            PostTradable: v.PostTradable,
            PostTotal:    v.PostTotal,
        }
        ll[i] = vv
    }
    return &RecordVol {
        Code: code,
        Data: ll,
    }, nil
}
func convertStkXDataToRecordXDR(code string, list []*gostk.XData) (*RecordXDR, error) {
    if len(list) <= 0 {
        return nil, ErrorInvalidParameter;
    }
    var ll = make([]*XData, len(list))
    for i, x := range list {
        xx := &XData {
            Date:          utils.DateTimeToDecimalNum(x.Date),
            AllotVolume:   x.AllotVolume,
            AllotPrice:    x.AllotPrice,
            BonusVolume:   x.BonusVolume,
            BonusPrice:    x.BonusPrice,
        }
        ll[i] = xx
    }
    return &RecordXDR {
        Code: code,
        Data: ll,
    }, nil
}

func convertRecordVolToStkVDataList(vol *RecordVol, from, to time.Time) ([]*gostk.VData) {
    var list = make([]*gostk.VData, len(vol.Data));
    var i = 0;
    for _, v := range vol.Data {
        t := utils.DecimalNumToDateTime(v.GetDate());
        if  t.Before(from) {
            continue;
        } else if t.After(to) {
            break;
        }
        vv := &gostk.VData {
            Date:           t,
            PreTradable:    v.GetPreTradable(),
            PreTotal:       v.GetPreTotal(),
            PostTradable:   v.GetPostTradable(),
            PostTotal:      v.GetPostTotal(),
        }

        list[i] = vv;
        i += 1
    }
    return list[:i];
}

func convertRecordXDRToStkXDataList(xdr *RecordXDR, from, to time.Time) ([]*gostk.XData) {
    var list = make([]*gostk.XData, len(xdr.Data));
    var i = 0;
    for _, x := range xdr.Data {
        t := utils.DecimalNumToDateTime(x.GetDate());
        if  t.Before(from) {
            continue;
        } else if t.After(to) {
            break;
        }
        xx := &gostk.XData {
            Date:          t,
            AllotVolume:   x.GetAllotVolume(),
            AllotPrice:    x.GetAllotPrice(),
            BonusVolume:   x.GetBonusVolume(),
            BonusPrice:    x.GetBonusPrice(),
        }
        list[i] = xx;
        i += 1
    }
    return list[:i];
}

func mergeRecordVol(to *RecordVol, from *RecordVol) *RecordVol {
    if to == nil && from == nil {
        return nil;
    }
    if to == nil {
        return from
    } else if from == nil {
        return to
    }
    list := make([]*VData, 0)
    var len_to int = len(to.Data)
    var len_from int = len(from.Data)
    var i int = 0
    var j int = 0
    for ;i < len_to && j < len_from; {
        if to.Data[i].Date == from.Data[j].Date { 
            list = append(list, to.Data[i])
            i += 1
            j += 1
        } else if to.Data[i].Date < from.Data[j].Date {
            list = append(list, to.Data[i])
            i += 1
        } else if to.Data[i].Date > from.Data[j].Date {
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
    vol := &RecordVol {
        Data: list,
    }
    return vol;
}


func mergeRecordXDR(to *RecordXDR, from *RecordXDR) *RecordXDR {
    if to == nil && from == nil {
        return nil;
    }
    if to == nil {
        return from
    } else if from == nil {
        return to
    }
    list := make([]*XData, 0)
    var len_to int = len(to.Data)
    var len_from int = len(from.Data)
    var i int = 0
    var j int = 0
    for ;i < len_to && j < len_from; {
        if to.Data[i].Date == from.Data[j].Date { 
            list = append(list, to.Data[i])
            i += 1
            j += 1
        } else if to.Data[i].Date < from.Data[j].Date {
            list = append(list, to.Data[i])
            i += 1
        } else if to.Data[i].Date > from.Data[j].Date {
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
    xdr := &RecordXDR {
        Data: list,
    }
    return xdr;
}

func getRecordVol(code string) (*RecordVol, error) {
    key := fmt.Sprintf("vol_%s", code); 
    bin, err := dbBase.Get([]byte(key), &opt.ReadOptions{})
    if err != nil {
        return nil, err;
    }
    vol := &RecordVol{};
    err = proto.Unmarshal(bin, vol)
    if err != nil {
        return nil, err;
    }
    //logger.Debugf("[LevelDB] getvol key: %v, len(vol.Data): %d", key, len(vol.Data))
    return vol, nil;
}

func putRecordVol(code string, vol *RecordVol) error {
    bin, err := proto.Marshal(vol);
    if err != nil {
        return err;
    }
    key := fmt.Sprintf("vol_%s", code); 
    err = dbBase.Put([]byte(key), bin, nil)
    if err != nil {
        return err;
    }
    //logger.Debugf("[LevelDB] putvol key: %v, len(vol.Data): %d", key, len(vol.Data))
    return nil
}

func getRecordXDR(code string) (*RecordXDR, error) {
    key := fmt.Sprintf("xdr_%s", code); 
    bin, err := dbBase.Get([]byte(key), &opt.ReadOptions{})
    if err != nil {
        return nil, err;
    }
    xdr := &RecordXDR{};
    err = proto.Unmarshal(bin, xdr)
    if err != nil {
        return nil, err;
    }
    //logger.Debugf("[LevelDB] getxdr key: %v, len(xdr.Data): %d", key, len(xdr.Data))
    return xdr, nil;
}

func putRecordXDR(code string, xdr *RecordXDR) error {
    bin, err := proto.Marshal(xdr);
    if err != nil {
        return err;
    }
    key := fmt.Sprintf("xdr_%s", code); 
    err = dbBase.Put([]byte(key), bin, nil)
    if err != nil {
        return err;
    }
    //logger.Debugf("[LevelDB] putxdr key: %v, len(xdr.Data): %d", key, len(xdr.Data))
    return nil
}

func GetVolumeList(code string, from, to time.Time) ([]*gostk.VData, error) {
    if from.After(to) {
        return nil, ErrorInvalidParameter;
    }
    if !gostk.IsStockCodeValid(code) {
        return nil, ErrorInvalidParameter;
    }
    vol, err := getRecordVol(code)
    if err != nil {
        return nil, err
    }
    ll := convertRecordVolToStkVDataList(vol, from, to)
    logger.Debugf("[LevelDB] GetVolumeList code(%s), date(%s - %s), size(%d)", 
                                            code, 
                                            from.Format("20060102"),
                                            to.Format("20060102"),
                                            len(ll))
    return ll, err
}

func PutVolumeList(code string, list []*gostk.VData) error {
    if len(list) == 0 || !gostk.IsStockCodeValid(code) {
        return ErrorInvalidParameter;
    }
    sort.SliceStable(list, func(i, j int) bool {
        return list[i].Date.Before(list[j].Date)
    })
    var err error = nil

    from, err := getRecordVol(code)
    if err != nil && err != errors.ErrNotFound {
        return err
    }
    to, err := convertStkVDataToRecordVol(code, list)
    if err != nil {
        return err;
    }
    newVol := mergeRecordVol(to, from)
    err = putRecordVol(code, newVol)
    from_len := 0;
    if from != nil {
        from_len = len(from.Data)
    }
    logger.Debugf("[LevelDB] PutVolumeList code(%s), size(%d, %d, %d)", 
                                                            code, 
                                                            from_len,
                                                            len(to.Data),
                                                            len(newVol.Data))
    return err;
}
func GetXDRList(code string, from, to time.Time) ([]*gostk.XData, error) {
    if from.After(to) {
        return nil, ErrorInvalidParameter;
    }
    if !gostk.IsStockCodeValid(code) {
        return nil, ErrorInvalidParameter;
    }
    xdr, err := getRecordXDR(code)
    if err != nil {
        return nil, err
    }
    ll := convertRecordXDRToStkXDataList(xdr, from, to)
    logger.Debugf("[LevelDB] GetXDRList code(%s), date(%s - %s), size(%d)", 
                                            code, 
                                            from.Format("20060102"),
                                            to.Format("20060102"),
                                            len(ll))
    return ll, err
}

func PutXDRList(code string, list []*gostk.XData) error {
    if len(list) == 0 || !gostk.IsStockCodeValid(code) {
        return ErrorInvalidParameter;
    }
    sort.SliceStable(list, func(i, j int) bool {
        return list[i].Date.Before(list[j].Date)
    })
    var err error = nil
    from, err := getRecordXDR(code)
    if err != nil && err != errors.ErrNotFound {
        return err
    }
    to, err := convertStkXDataToRecordXDR(code, list)
    if err != nil {
        return err;
    }
    newXdr := mergeRecordXDR(to, from)
    err = putRecordXDR(code, newXdr)
    from_len := 0;
    if from != nil {
        from_len = len(from.Data)
    }
    logger.Debugf("[LevelDB] PutXDRList code(%s), size(%d, %d, %d)", 
                                                            code, 
                                                            from_len,
                                                            len(to.Data),
                                                            len(newXdr.Data))
    return err;
}
