package gostk


/*
 前复权
 从最近K线前开始计算，由K线日起开始往前搜索除权信息， 如果K线日 < (早于)除权日， 计算前复权
 kdata 和 xdr 需要确保为增序
**/
func ForwardAdjuste(kdataList []*KData, xdrList []*XData) []*KData {
    kdata_len := len(kdataList)
    xdr_len := len(xdrList)
    l := make([]*KData, kdata_len);

    for i := kdata_len - 1; i >= 0; i-- {
        k := kdataList[i];
        k1 := &KData{
            Time: k.Time,
            Open: k.Open,
            Close: k.Close,
            High: k.High,
            Low: k.Low,
            Volume: k.Volume,
            Turnover: k.Turnover,
        }
        for j := xdr_len - 1; j >= 0 ; j-- {
            x := xdrList[j]
            if k.Time.Before(x.Date) {
                if x.Type == XData_EXP {
                    k1.Open  = k1.Open / x.BonusVolume;
                    k1.Close = k1.Close / x.BonusVolume;
                    k1.High  = k1.High / x.BonusVolume;
                    k1.Low   = k1.Low / x.BonusVolume;
                } else {
                    ratio := float32(x.BonusVolume / 10.0 + 1);
                    price := float32(x.BonusPrice / 10.0)
                    k1.Open  = k1.Open / ratio - price;
                    k1.Close = k1.Close / ratio - price;
                    k1.High  = k1.High / ratio - price;
                    k1.Low   = k1.Low / ratio - price;
                }
            }
        }
        l[i] = k1;
    }
    return l;
}

/* 
  后复权 
  从最早的K线数据开始，由K线日起开始往前搜索除权信息，如果K线日期>= (晚于或等于) 除权当日，计算后复权
  kdata 和 xdr 需要确保为增序
**/
func BackwardAdjust(kdataList []*KData, xdrList []*XData) []*KData {
	kdata_len := len(kdataList)
	xdr_len := len(xdrList)
	l := make([]*KData, kdata_len);
	for i := 0; i < kdata_len; i++ {
		k := kdataList[i];
		k1 := &KData{
			Time: k.Time,
			Open: k.Open,
			Close: k.Close,
			High: k.High,
			Low: k.Low,
			Volume: k.Volume,
			Turnover: k.Turnover,
		}
		for j := xdr_len - 1; j >= 0; j-- {
			x := xdrList[j]
			if k.Time.After(x.Date) || k.Time.Equal(x.Date) {
				if x.Type == XData_EXP {
					k1.Open  = k1.Open * x.BonusVolume;
					k1.Close = k1.Close * x.BonusVolume;
					k1.High  = k1.High * x.BonusVolume;
					k1.Low   = k1.Low * x.BonusVolume;
				} else {
					ratio := float32(x.BonusVolume / 10.0 + 1);
					price := float32(x.BonusPrice / 10.0)
					k1.Open  = k1.Open * ratio + price;
					k1.Close = k1.Close * ratio + price;
					k1.High  = k1.High * ratio + price;
					k1.Low   = k1.Low * ratio + price;
				}
			}
		}
		l[i] = k1;
	}
	return l;
}
