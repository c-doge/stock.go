package gostk

import (
	"strings"
)

type MarketType int32 

const (
	MT_UKN               MarketType = 0
	MT_SHSE              MarketType = 1
	MT_SZSE              MarketType = 2
)

type StockType int32

const (
	STK_UKN              StockType = 0

    STK_SHSE_IDX         StockType = 1
    STK_SHSE_A           StockType = 2
    STK_SHSE_B           StockType = 3
    STK_SHSE_STAR        StockType = 4 // 科创
    
    STK_SZSE_IDX         StockType = 1001
    STK_SZSE_A           StockType = 1002
    STK_SZSE_B           StockType = 1003
    STK_SZSE_SME         StockType = 1004
    STK_SZSE_GEM         StockType = 1005
)

func IsIndex(t StockType) bool {
	if t == STK_SHSE_IDX || t == STK_SZSE_IDX {
		return true;
	}
	return false;
}

func IsStockCodeValid(code string) bool {
    if strings.Contains(code,  "sh60")  ||    // 沪市主板A
        strings.Contains(code, "sh90")  ||    // 沪市主板B
        strings.Contains(code, "sh00")  ||    // 沪市指数
        strings.Contains(code, "sh999") ||    // 沪市指数
        strings.Contains(code, "sh688") ||    // 沪市科创板
        strings.Contains(code, "sh204") ||    // 沪市逆回国
        strings.Contains(code, "sz000") ||    // 深市主板A
        strings.Contains(code, "sz20")  ||    // 深市主板B
        strings.Contains(code, "sz002") ||    // 深市中小板
        strings.Contains(code, "sz300") ||    // 深市创业板
        strings.Contains(code, "sz399") ||    // 深市板块
        strings.Contains(code, "sz1318") {    // 深市逆回购
        return true;
    }
    return false;
}

func GetStockTypeFromCode(code string) StockType {
	if code[1] == 'h' {
		if strings.HasPrefix(code, "sh00") || strings.HasPrefix(code, "sh999") {
			return STK_SHSE_IDX;
		} else if strings.HasPrefix(code, "sh60") {
			return STK_SHSE_A;
		} else if strings.HasPrefix(code, "sh90") {
			return STK_SHSE_B;
		} else if strings.HasPrefix(code, "sh688") {
			return STK_SHSE_STAR;
		} 
	} else if code[1] == 'z' {
		if strings.HasPrefix(code, "sz395") || strings.HasPrefix(code, "sz399") {
			return STK_SZSE_IDX;
		} else if strings.HasPrefix(code, "sz000") {
			return STK_SZSE_A;
		} else if strings.HasPrefix(code, "sz002") {
			return STK_SZSE_SME;
		} else if strings.HasPrefix(code, "sz20") {
			return STK_SZSE_B;
		} else if strings.HasPrefix(code, "sz300") {
			return STK_SZSE_GEM;
		} 
	}
	return STK_UKN
}

func GetMarketTypeFromStockType(t StockType) MarketType {
	if t >= STK_SHSE_IDX && t <= STK_SHSE_STAR {
		return MT_SHSE
	} else if t >= STK_SZSE_IDX && t <= STK_SZSE_GEM {
		return MT_SZSE
	}
	return MT_UKN;
}

func GetStockTypeLabel(t StockType) string {
	if t == STK_SHSE_IDX {
		return "沪市指数"
	} else if t == STK_SHSE_A {
		return "沪市主板A股"
	} else if t == STK_SHSE_B {
		return "沪市主板B股"
	} else if t == STK_SHSE_STAR {
		return "科创板"
	} else if t == STK_SZSE_IDX {
		return "深市指数"
	} else if t == STK_SZSE_A {
		return "深市主板A股"
	} else if t == STK_SZSE_B {
		return "深市主板B股"
	} else if t == STK_SZSE_SME {
		return "中小板"
	} else if t == STK_SZSE_GEM {
		return "创业板"
	}
	return "未知类型"
}