package leveldb

import (
	"fmt"
	"time"
	"math"
	"testing"
	"math/rand"
    "github.com/c-doge/stock.go/gostk"
    "github.com/c-doge/stock.go/base/utils"
)
func equalVData(d1, d2 *gostk.VData) bool {
    if !d1.Date.Equal(d2.Date) {
        return false
    } else if math.Abs(d1.PreTradable - d2.PreTradable) > 0.00001 {
        return false
    } else if math.Abs(d1.PostTradable - d2.PostTradable) > 0.00001 {
        return false
    } else if math.Abs(d1.PreTotal - d2.PreTotal) > 0.00001 {
        return false
    } else if math.Abs(d1.PostTotal - d2.PostTotal) > 0.00001 {
        return false
    }
    return true;
}
func equalXData(d1, d2 *gostk.XData) bool {
    if !d1.Date.Equal(d2.Date) {
        return false
    } else if d1.Type != d2.Type {
        return false
    } else if abs32(d1.AllotVolume, d2.AllotVolume) > 0.00001 {
        return false
    } else if abs32(d1.AllotPrice, d2.AllotPrice) > 0.00001 {
        return false
    } else if abs32(d1.BonusVolume, d2.BonusVolume) > 0.00001 {
        return false
    } else if abs32(d1.BonusPrice, d2.BonusPrice) > 0.00001 {
        return false
    }
    return true;
}

func newStkVData(t time.Time) *gostk.VData {
	pre := rand.Float64() * 1000000
	post := pre + rand.Float64() * 1000
    v := &gostk.VData {
    	Date:          t,
        PreTradable:   pre,
        PreTotal:      pre + rand.Float64() * 1000,
        PostTradable:  post,
        PostTotal:	   post + rand.Float64() * 1000,	
    };
    return v;
}

func newStkXData(t time.Time) *gostk.XData {
	newType := func() gostk.XDataType {
		if rand.Float32() > 0.5 {
			return gostk.XData_EXP
		}
		return gostk.XData_XDR
	}	
    x := &gostk.XData {
    	Date:          t,
    	Type:          newType(),
        AllotVolume:   rand.Float32() * 100,
        AllotPrice:    rand.Float32() * 10,
        BonusVolume:   rand.Float32() * 100,
        BonusPrice:	   rand.Float32() * 10,	
    };
    return x;
}

func checkVolumeList(l1, l2 []*gostk.VData) error {
	if len(l1) != len(l2) {
		return fmt.Errorf("The two VData list are not equal in length")
	}
	for i := 0; i < len(l1); i++ {
		if !equalVData(l1[i], l2[i]) {
			return fmt.Errorf("VData[%s] is not equal to VData[%s], index(%d)", 
							l1[i].Date.Format("2006-01-02"), 
							l2[i].Date.Format("2006-01-02"),
							i)
		}
	}
	return nil
}

func checkXDRList(l1, l2 []*gostk.XData) error {
	if len(l1) != len(l2) {
		return fmt.Errorf("The two XData list are not equal in length")
	}
	for i := 0; i < len(l1); i++ {
		if !equalXData(l1[i], l2[i]) {
			return fmt.Errorf("XData[%s] is not equal to XData[%s], index(%d)", 
							l1[i].Date.Format("2006-01-02"), 
							l2[i].Date.Format("2006-01-02"),
							i)
		}
	}
	return nil
}

func Test_volPutAndGet(t *testing.T) {
	code := "sh000001"

	date := utils.DecimalNumToDateTime(20160205)
	l1 := make([]*gostk.VData, 0)
	for i := 0; i < 10; i++ {
		v := newStkVData(date)
		l1 = append(l1, v)
		date = date.AddDate(0, 0 ,1)
	}

	err := PutVolumeList(code, l1)
	if err != nil {
		t.Errorf("volPutAndGet, PutVolumeList fail, Error:%v", err);
		return
	}
	l_result, err := GetVolumeList(code, utils.DecimalNumToDateTime(20160205), utils.DecimalNumToDateTime(20160220))
	if err != nil {
		t.Errorf("volPutAndGet, GetVolumeList fail, Error:%v", err);
		return
	}
	err = checkVolumeList(l1, l_result) 
	if err != nil {
		t.Errorf("volPutAndGet, checkVolumeList %v", err);
		return
	}

	l2 := make([]*gostk.VData, 0)
	for i := 0; i < 5; i++ {    //20160201 - 20160210
		v := newStkVData(date)
		l2 = append(l2, v)
		date = date.AddDate(0, 0 ,1)
	}

	err = PutVolumeList(code, l2)
	if err != nil {
		t.Errorf("volPutAndGet.Append, PutVolumeList fail, Error:%v", err);
		return
	}
	l_result, err = GetVolumeList(code, utils.DecimalNumToDateTime(20160205), utils.DecimalNumToDateTime(20160310))
	if err != nil {
		t.Errorf("volPutAndGet.Append, GetVolumeList fail, Error:%v", err);
		return
	}
	l3 := append(l1, l2...)
	err = checkVolumeList(l3, l_result) 
	if err != nil {
		t.Errorf("volPutAndGet.Append, checkVolumeList %v", err);
		return
	}
}

func Test_xdrPutAndGet(t *testing.T) {
	code := "sh000001"

	date := utils.DecimalNumToDateTime(20160205)
	l1 := make([]*gostk.XData, 0)
	for i := 0; i < 10; i++ {
		v := newStkXData(date)
		l1 = append(l1, v)
		date = date.AddDate(0, 0 ,1)
	}

	err := PutXDRList(code, l1)
	if err != nil {
		t.Errorf("xdrPutAndGet, PutXDRList fail, Error:%v", err);
		return
	}
	l_result, err := GetXDRList(code, utils.DecimalNumToDateTime(20160205), utils.DecimalNumToDateTime(20160220))
	if err != nil {
		t.Errorf("xdrPutAndGet, GetXDRList fail, Error:%v", err);
		return
	}
	err = checkXDRList(l1, l_result) 
	if err != nil {
		t.Errorf("xdrPutAndGet, checkXDRList %v", err);
		return
	}

	l2 := make([]*gostk.XData, 0)
	for i := 0; i < 5; i++ {
		v := newStkXData(date)
		l2 = append(l2, v)
		date = date.AddDate(0, 0 ,1)
	}

	err = PutXDRList(code, l2)
	if err != nil {
		t.Errorf("xdrPutAndGet.Append, PutXDRList fail, Error:%v", err);
		return
	}
	l_result, err = GetXDRList(code, utils.DecimalNumToDateTime(20160205), utils.DecimalNumToDateTime(20160310))
	if err != nil {
		t.Errorf("xdrPutAndGet.Append, GetXDRList fail, Error:%v", err);
		return
	}
	l3 := append(l1, l2...)
	err = checkXDRList(l3, l_result) 
	if err != nil {
		t.Errorf("xdrPutAndGet.Append, checkXDRList %v", err);
		return
	}
}