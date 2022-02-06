package db

import (
	"time"
	"github.com/c-doge/stock.go/gostk"	
	"github.com/c-doge/stock.go/db/leveldb"
	
)

type leveldbHelper struct {
	dbPath string
}

func (helper *leveldbHelper) Start(dbPath string) error  {
	dbPath = dbPath
	return leveldb.Start(dbPath)
}
func (helper *leveldbHelper) Stop() {
	leveldb.Stop();
}
func (helper *leveldbHelper) PutLday(code string, list []*gostk.KData) error {
	return leveldb.PutLday(code, list)
}
func (helper *leveldbHelper) GetLday(code string, from, to time.Time) ([]*gostk.KData, error) {
	return leveldb.GetLday(code, from, to)
}

