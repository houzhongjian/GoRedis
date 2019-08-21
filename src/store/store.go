package store

import (
	"log"

	"github.com/houzhongjian/GoRedis/src/conf"
	"github.com/syndtr/goleveldb/leveldb"
)

//todo 存储以后采用rocksdb，目前采用map方便开发测试.
var Data = make(map[string]interface{})

//StoreEngine .
type StoreEngine struct {
	db        *leveldb.DB
	storePath string
}

func New() *StoreEngine {
	return newStoreEngine()
}

//newStoreEngine .
func newStoreEngine() *StoreEngine {
	var storePath = "./store/"
	if conf.IsExist("storepath") {
		storePath = conf.GetString("storepath")
	}
	eng := &StoreEngine{
		storePath: storePath,
	}

	eng.newDB()
	return eng
}

func (s *StoreEngine) newDB() {
	db, err := leveldb.OpenFile(s.storePath, nil)
	if err != nil {
		log.Panicf("%+v\n", err)
		return
	}
	s.db = db
}

//Insert 记录一条数据.
func (s *StoreEngine) Insert(key, value string) error {
	return s.db.Put([]byte(key), []byte(value), nil)
}

//Query查询一条数据.
func (s *StoreEngine) Query(key string) (b []byte, err error) {
	b, err = s.db.Get([]byte(key), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return b, err
	}
	return b, nil
}
