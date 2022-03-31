package gostk

import (
    "testing"
    "math"
    "github.com/c-doge/stock.go/base/utils"
    "github.com/c-doge/stock.go/base/logger"
)


func abs32(f1, f2 float32) float32 {
    var v float32 = f1 - f2;
    if v > 0 {
        return v;
    }
    return 0 - v;
}

func equalKData(d1, d2 *KData) bool {
    if !d1.Time.Equal(d2.Time) {
        logger.Debugf("time: %s,%s", d1.Time.Format("2006-01-02"), d2.Time.Format("2006-01-02"))
        return false
    } else if abs32(d1.Open, d2.Open) > 0.006 {
        logger.Debugf("%s, open: %f,%f", d1.Time.Format("2006-01-02"), d1.Open, d2.Open)
        return false
    } else if abs32(d1.Close, d2.Close) > 0.006 {
        logger.Debugf("%s, close: %f,%f", d1.Time.Format("2006-01-02"), d1.Close, d2.Close)
        return false
    } else if abs32(d1.High, d2.High) > 0.006 {
        logger.Debugf("%s, high: %f,%f", d1.Time.Format("2006-01-02"), d1.High, d2.High)
        return false
    } else if abs32(d1.Low, d2.Low) > 0.006 {
        logger.Debugf("%s, low: %f,%f", d1.Time.Format("2006-01-02"), d1.Low, d2.Low)
        return false
    } else if math.Abs(d1.Volume - d2.Volume) > 0.006 {
        logger.Debugf("%s, Volume, %f-%f", d1.Time.Format("2006-01-02"), d1.Volume, d2.Volume)
        return false
    } else if math.Abs(d1.Turnover - d2.Turnover) > 0.006 {
        logger.Debugf("%s, Turnover: %f-%f", d1.Time.Format("2006-01-02"), d1.Turnover, d2.Turnover)
        return false
    }
    return true;
}

func equalKDataList(k1 []*KData, k2 []*KData) bool {
	if len(k1) != len(k2) {
		return false;
	}
	l := len(k1)
	for i := 0; i < l; i++ {
		if !equalKData(k1[i], k2[i]) {
			return false;
		}
	}
	return true;
}

func getSH600381Xdr() []*XData {
	xdrList := []*XData {
		&XData{Date: utils.DecimalNumToDateTime(20020705), Type: XData_XDR,  BonusPrice: 0.650, AllotPrice: 0, BonusVolume: 0,        AllotVolume: 0,},
		&XData{Date: utils.DecimalNumToDateTime(20030521), Type: XData_XDR,  BonusPrice: 0.680, AllotPrice: 0, BonusVolume: 0,        AllotVolume: 0,},
		&XData{Date: utils.DecimalNumToDateTime(20040719), Type: XData_XDR,  BonusPrice: 0,     AllotPrice: 0, BonusVolume: 10.0,     AllotVolume: 0,},
		&XData{Date: utils.DecimalNumToDateTime(20050512), Type: XData_XDR,  BonusPrice: 0,     AllotPrice: 0, BonusVolume: 3.0,      AllotVolume: 0,},
		&XData{Date: utils.DecimalNumToDateTime(20061127), Type: XData_XDR,  BonusPrice: 0,     AllotPrice: 0, BonusVolume: 2.24,     AllotVolume: 0,},
		&XData{Date: utils.DecimalNumToDateTime(20110825), Type: XData_XDR,  BonusPrice: 0,     AllotPrice: 0, BonusVolume: 5.0,      AllotVolume: 0,},
		&XData{Date: utils.DecimalNumToDateTime(20120515), Type: XData_XDR,  BonusPrice: 0,     AllotPrice: 0, BonusVolume: 7.0,      AllotVolume: 0,},
		&XData{Date: utils.DecimalNumToDateTime(20140627), Type: XData_EXP,  BonusPrice: 0,     AllotPrice: 0, BonusVolume: 0.124185, AllotVolume: 0,},
	}
	return xdrList
}

func Test_Adjust(t *testing.T) {
	xdr := getSH600381Xdr();
	rawLday, err := ReadLdayCsv("../test/csv/sh600381-raw.txt");
	if err != nil {
		t.Errorf("TestAdjust, Read sh600381-raw.txt Error: %v", err);
		return;
	}
	fLDay, err := ReadLdayCsv("../test/csv/sh600381-forward.txt");
	if err != nil {
		t.Errorf("TestAdjust, Read sh600381-forward.txt Error: %v", err);
		return;
	}
	bLday, err := ReadLdayCsv("../test/csv/sh600381-backward.txt");
	if err != nil {
		t.Errorf("TestAdjust, Read sh600381-backward.txt Error: %v", err);
		return;
	}

	l_forward := ForwardAdjuste(rawLday, xdr)
	if !equalKDataList(l_forward, fLDay) {
		t.Error("TestAdjust, ForwardAdjuste Fail");
		return;
	}

	l_backward := BackwardAdjust(rawLday, xdr)
	if !equalKDataList(l_backward, bLday) {
		t.Error("TestAdjust, BackwardAdjust Fail");
		return;
	}
	return;
}

func TestMain(m *testing.M) {
    logger.New("Debug", "", "stock.go/gostk")
    logger.Info("stock.go gostk test start >>>")
    m.Run()
    logger.Info("stock.go gostk test stop >>>")
}