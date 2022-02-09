package leveldb


import (
	"path"
	"errors"
//    "github.com/golang/protobuf/proto"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"

	// "github.com/c-doge/stock.go/gostk"
	// "github.com/c-doge/stock.go/base/utils"
	"github.com/c-doge/stock.go/base/logger"
)

const (
	dbFileLday = "lday.db" 
	dbFileBase = "base.db"

)

var dbLday *leveldb.DB = nil;
var dbBase *leveldb.DB = nil;

var ErrorNotResult         error = errors.New("Not Result");
var ErrorInvalidParameter  error = errors.New("Invalid Parameter");

func Start(dbPath string) error {
	logger.Infof("[LevelDB] Initialize");
	var err error = nil
	o := &opt.Options{
		Compression: opt.SnappyCompression,
	}
	// BaseDB
	dbPathBase := path.Join(dbPath, dbFileBase)
	dbBase, err = leveldb.OpenFile(dbPathBase, o)
	if err != nil {
		logger.Fatalf("[LevelDB] open base db fail, err: %v", err);
		return err;
	}
	err = checkVolVersion() 
	if err != nil {
		return err
	}

	// LdayDB
	dbPathLday := path.Join(dbPath, dbFileLday)
	dbLday, err = leveldb.OpenFile(dbPathLday, o)
	if err != nil {
		logger.Fatalf("[LevelDB] open lday db fail, err: %v", err);
		return err;
	}
	err = checkLdayVersion() 
	if err != nil {
		return err
	}
	return nil;
}

func Stop() {
	logger.Infof("[LevelDB] Close ")
	if dbLday != nil {
		dbLday.Close();
		dbLday = nil;
	}
	if dbBase != nil {
		dbBase.Close();
		dbBase = nil;
	}
}




