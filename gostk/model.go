package gostk

import (
    "fmt"
    "time"
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

    result := fmt.Sprintf("%v, %v, %10.3f, %10.3f, %10.3f, %10.3f, %15.3f, %20.3f",
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

// 除权/除息 数据
type XData struct {
    Date                    time.Time     //"时间戳" 
    // 配股
    AllotVolume             float32       // 每十股配股数
    AllotPrice              float32       // 配股价格
    // 分红/送股
    BonusVolume             float32       // 每十股送股数
    BonusPrice              float32       // 每十股分红数
}

func (data XData) String () string {

    result := fmt.Sprintf("%v,  %12.3f,  %12.3f,  %12.3f,  %12.3f",
        data.Date.Format("2006-01-02"),
        data.AllotPrice,
        data.AllotVolume,
        data.BonusPrice,
        data.BonusVolume);

    return result;
}

// 股本信息
type VData struct {
    Date               time.Time      //"时间戳"
    Tradable           float64        // 流通股  x10000
    Total              float64        // 总股本  x10000
}

func (data VData) String () string {
    result := fmt.Sprintf("%v, %f, %f",
        data.Date.Format("2006-01-02"),
        data.Tradable,
        data.Total);
    return result;
}