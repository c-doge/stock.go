package leveldb

import (
    "time"
    "math"
    "testing"
    "math/rand"

    "github.com/c-doge/stock.go/gostk"
    "github.com/c-doge/stock.go/base/utils"
    "github.com/c-doge/stock.go/base/logger"
)

func equalLeveldbataKData(d1, d2 *KData) bool {
    if d1.GetTime() != d2.GetTime() {
        return false
    } else if abs32(d1.GetOpen(), d2.GetOpen()) > 0.00001 {
        return false
    } else if abs32(d1.GetClose(), d2.GetClose()) > 0.00001 {
        return false
    } else if abs32(d1.GetHigh(), d2.GetHigh()) > 0.00001 {
        return false
    } else if abs32(d1.GetLow(), d2.GetLow()) > 0.00001 {
        return false
    } else if math.Abs(d1.GetVolume() - d2.GetVolume()) > 0.00001 {
        return false
    } else if math.Abs(d1.GetTurnover() - d2.GetTurnover()) > 0.00001 {
        return false
    }
    return true;
}

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

func newStkKDataList(year, size int) []*gostk.KData {
    date := uint32(year*10000+101)
    t := utils.DecimalNumToDateTime(date)
    ll := make([]*gostk.KData, 0)
    for i:= 0; i < size; i+=1 {
        kdata:= newStkKData(t);
        t = t.AddDate(0, 0, 1)
        ll = append(ll,  kdata)
    }
    return ll
}

func newLeveldbKData(t time.Time) *KData {
    kdata := newStkKData(t);
    return &KData {
        Time:     utils.DateTimeToDecimalNum(kdata.Time),
        Open:     kdata.Open,
        Close:    kdata.Close,
        High:     kdata.High,
        Low:      kdata.Low,
        Volume:   kdata.Volume,
        Turnover: kdata.Turnover,
    }
}

func newRecordLday(year, size, start, step int) *RecordLday {
    date := uint32(year*10000+start)
    t := utils.DecimalNumToDateTime(date)
    ll := make([]*KData, 0)
    for i:= 0; i < size; i+=1 {
        kdata:= newLeveldbKData(t);
        t = t.AddDate(0, 0, step)
        ll = append(ll,  kdata)
    }
    lday := &RecordLday{}
    lday.Data = ll;
    return lday
}

func dumpRecordLday(lday *RecordLday) {
    for i, v := range lday.Data {
        logger.Debugf("i:%d, %v", i, v);
    }
}

func Test_ldayConvert(t *testing.T) {

    l := newStkKDataList(2016, 302)
    lday, err := convertStkKDataToRecordLday("sh601888", l)
    if err != nil {
        t.Errorf("convert kdata to dayline fail, err: %v\n", err)
    }
    if len(l) != len(lday.GetData()) && len(l) != len(lday.Data) {
        t.Errorf("length of KData not equal dayline\n")
        return;
    }
    ll := convertRecordLdayToStkKDataList(lday, utils.EarlyDate, utils.FutureDate);

    if len(l) != len(ll) {
        t.Errorf("length of KData not equal after convert, %d, %d\n", len(l), len(ll));
        return;
    }
    for i, v := range l {
        if !v.Equal(*ll[i]) {
            t.Errorf("kdata not equal i = %d, \n", i);
            return;
        }
    }
}

func Test_mergeRecordLday(t *testing.T) {
    l1_1 := newRecordLday(2016, 10, 101, 2)
    l1_2 := newRecordLday(2016, 15, 102, 2)
    ll1 := mergeRecordLday(l1_1, l1_2)

    if 25 != len(ll1.Data) {
        t.Errorf("Case 1, size of ll1(%d) not equal to 25", len(ll1.Data));
    }
    for i := 0; i < 10; i++ {
        if !equalLeveldbataKData(l1_1.Data[i], ll1.Data[i*2]) {
            t.Errorf("Case 1, ll1[%d]  not equal l1_1[%d], \n", i*2, i);
        } else if !equalLeveldbataKData(l1_2.Data[i], ll1.Data[i*2 + 1]) {
            t.Errorf("Case 1, ll1[%d]  not equal l1_2[%d], \n", i*2+1, i);
        }
    }

    for i := 0; i < 5; i++ {
        if !equalLeveldbataKData(l1_2.Data[10+i], ll1.Data[20 + i]) {
            t.Errorf("Case 1, ll1[%d]  not equal l1_2[%d], \n", 20+i, 10+i);
        }
    }
    
    l2_1 := newRecordLday(2016, 15, 101, 2)
    l2_2 := newRecordLday(2016, 10, 102, 2)
    ll2 := mergeRecordLday(l2_1, l2_2)

    if 25 != len(ll2.Data) {
        t.Errorf("Case 2, size of ll2(%d) not equal to 25", len(ll2.Data));
    }
    for i := 0; i < 10; i++ {
        if !equalLeveldbataKData(l2_1.Data[i], ll2.Data[i*2]) {
            t.Errorf("Case 2, ll2[%d]  not equal l2_1[%d], \n", i*2, i);
        } else if !equalLeveldbataKData(l2_2.Data[i], ll2.Data[i*2 + 1]) {
            t.Errorf("Case 2, ll2[%d]  not equal l2_2[%d], \n", i*2+1, i);
        }
    }

    for i := 0; i < 5; i++ {
        if !equalLeveldbataKData(l2_1.Data[10+i], ll2.Data[20 + i]) {
            t.Errorf("Case 2, ll2[%d]  not equal l2_1[%d], \n", 20+i, 10+i);
        }
    }
    return
}

func Test_ldayPutAndGet(t *testing.T) {
    var err error = nil;
    var l_result []*gostk.KData = nil;
    var code = "sh601888"
    l_2016 := newStkKDataList(2016, 301)
    l_2017 := newStkKDataList(2017, 302)
    l_2018 := newStkKDataList(2018, 303)
    l := append(l_2016, l_2017...)
    l = append(l, l_2018...)
    err = PutLday(code, l);
    if err != nil {
        t.Errorf("PutDayline fail, err: %v\n", err)
        return
    }
    t_start := utils.DecimalNumToDateTime(20160101)
    t_end := utils.DecimalNumToDateTime(20171231)

    l_result, err = GetLday(code, t_start, t_end)
    if err != nil {
        t.Errorf("GetDayLine fail, err: %v\n", err)
        return;
    }

    if len(l_result) != len(l_2016) + len(l_2017) {
        t.Errorf("size of l_result not equal to l_2016 + l_2017");
        return
    }
    t_start = utils.DecimalNumToDateTime(20160501)
    t_end = utils.DecimalNumToDateTime(20180430)
    l_result, err = GetLday(code, t_start, t_end)
    if err != nil {
        t.Errorf("GetDayLine fail, err: %v\n", err)
        return;
    }

    ll := make([]*gostk.KData, 0)
    for _, v := range l {
        if v.Time.Before(t_start) || v.Time.After(t_end) {
            continue
        }
        ll = append(ll, v)
    }
    if len(l_result) != len(ll) {
        t.Errorf("size of l_result not equal to ll");
        return
    }
    for i, v := range l_result {
        if !v.Equal(*ll[i]) {
            t.Errorf("l_result's kdata not equal ll[%d], \n", i);
            return;
        }
    }
}
