
package gostk

import (
    "fmt"
    "time"
    "container/list"
//    "github.com/c-doge/stock.go/base/utils"
)

// K线
type KData struct {
    Time                    time.Time              //"时间戳"
    Open                    float32                //"开盘价"
    Close                   float32                //"收盘价"
    High                    float32                //"最高价"
    Low                     float32                //"最低价"
    Volume                  float64                //"成交量"
    Turnover                float64                //"成交额"
};

func (data KData) String () string {

    result := fmt.Sprintf("%s, %10.3f, %10.3f, %10.3f, %10.3f, %15.3f, %20.3f",
        data.Time.Format("2006-01-02 15:04"),
        data.Open,
        data.Close,
        data.High,
        data.Low,
        data.Volume,
        data.Turnover);

    return result;
}

func (data KData) Equal(kdata KData) bool {
    if data.Time.Equal(kdata.Time) &&
       data.Open == kdata.Open &&
       data.Close == kdata.Close &&
       data.High == kdata.High &&
       data.Low == kdata.Low &&
       data.Volume == kdata.Volume &&
       data.Turnover == kdata.Turnover {
        return true;
    }
    return false;
}

func (data KData) IsValid() bool{
    if data.Close <= 0 ||
        data.Open <= 0 ||
        data.High <= 0 ||
        data.Low <= 0 ||
        data.Volume <= 0 ||
        data.Turnover <= 0 {
        return false;
    }
    return true;
}

type KDataMap map[string]*list.List;

func NewKDataMap() KDataMap {
    return KDataMap{}
} 

func (m KDataMap) Insert(code string, value *KData) {
    l, ok := m[code];
    if !ok || l == nil {
        l = list.New();
    } 
    var pos *list.Element = nil;
    if l.Len() > 0 {
        for e := l.Front(); e != nil; e = e.Next() {
            v := e.Value.(*KData)
            if !value.Time.After(v.Time) {
                pos = e;
                break;
            }
        }
    }
    if pos == nil {
        l.PushBack(value)
    } else {
        l.InsertBefore(value, pos)
    }
    m[code] = l;
}
func (m KDataMap) HasKey(code string) bool {
    _, ok := m[code];
    return ok;
}
func (m KDataMap) Empty() bool {
    return len(m) == 0
}
func (m KDataMap) Size() int {
    return len(m)
}
func (m KDataMap) GetList(code string) *list.List {
    l, _ := m[code];
    return l;
}
func (m KDataMap) Get(code string) []*KData {
    l := m.GetList(code)
    if l == nil {
        return nil
    }
    ll := make([]*KData, 0, l.Len());
    for e := l.Front(); e != nil; e = e.Next() {
        v := e.Value.(*KData)
        ll = append(ll, v)
    }
    return ll
}
func (m KDataMap) ForEach(f func(code string, l *list.List)) {
    if f == nil {
        err := "KDataMap.ForEach fatal error without f";
        //logger.Fatalf(err)
        panic(err)
    }
    for code, _list := range m {
        f(code, _list)
    }
}

type XDataType uint32

const (
    XData_UKN              XDataType = 0
    XData_XDR              XDataType = 1
    XData_EXP              XDataType = 2
)

func (t XDataType) String() string {
    if t == XData_EXP {
        return "EXP"
    }
    return "XDR"
}

func XDataTypeFromString(s string) XDataType {
    if s == "EXP" {
        return XData_EXP;
    } else if s == "XDR" {
        return XData_XDR;
    }
    return XData_UKN;
}

// 除权/除息 数据
type XData struct {
    Date                    time.Time     //"时间戳" 
    Type                    XDataType     //类型 1 分红配股， 11      
    // 配股
    AllotVolume             float32       // 每十股配股数
    AllotPrice              float32       // 配股价格
    // 分红/送股
    BonusVolume             float32       // 每十股送股数
    BonusPrice              float32       // 每十股分红数
}

func (data XData) String () string {
    t := "除权除息"
    if data.Type == XData_EXP {
        t = "扩缩股  "
    }
    result := fmt.Sprintf("%v,  %s,    %12.6f,  %12.6f,  %12.6f,  %12.6f",
        data.Date.Format("2006-01-02"),
        t,
        data.AllotPrice,
        data.AllotVolume,
        data.BonusPrice,
        data.BonusVolume);

    return result;
}

// 股本信息
type VData struct {
    Date               time.Time      //"时间戳"
    PreTradable        float64        // 变动前 流通股  x10000
    PreTotal           float64        // 变动前 总股本  x10000
    PostTradable       float64        // 变动后 流通股  x10000
    PostTotal          float64        // 变动后 总股本  x10000
}

func (data VData) String () string {
    result := fmt.Sprintf("%v,     %0.6f,     %0.6f,     %0.6f,     %0.6f",
        data.Date.Format("2006-01-02"),
        data.PreTradable,
        data.PostTradable,
        data.PreTotal,
        data.PostTotal);
    return result;
}