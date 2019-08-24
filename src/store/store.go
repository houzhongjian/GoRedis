package store

import (
	"fmt"
	"log"
	"os"

	"github.com/houzhongjian/GoRedis/src/conf"
	"github.com/syndtr/goleveldb/leveldb"
)

//StoreEngine .
type StoreEngine struct {
	db        *leveldb.DB
	storePath string
}

func New() []*StoreEngine {
	return newStoreEngine()
}

//newStoreEngine .
func newStoreEngine() []*StoreEngine {
	var storePath = "./store/"
	if conf.IsExist("storepath") {
		storePath = conf.GetString("storepath")
	}

	//读取redis配置文件获取设置了多少个db.
	var databases = 16
	if conf.IsExist("databases") {
		databases = conf.GetInt("databases")
	}

	engList := []*StoreEngine{}
	//初始化db引擎.
	for i := 0; i < databases; i++ {
		eng := &StoreEngine{
			storePath: storePath,
		}

		//实例化数据库.
		eng.initDatabase(i)
		engList = append(engList, eng)
	}

	return engList
}

func (s *StoreEngine) newDB(selectdb int) {
	path := fmt.Sprintf("%s/db%d", s.storePath, selectdb)
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Panicf("%+v\n", err)
		return
	}
	s.db = db
}

func (s *StoreEngine) initDatabase(db int) {
	dbpath := fmt.Sprintf("%sdb%d", conf.GetString("storepath"), db)
	if err := os.MkdirAll(dbpath, os.ModePerm); err != nil {
		log.Printf("%+v\n", err)
		return
	}

	s.newDB(db)
}

//Insert 记录一条数据.
func (s *StoreEngine) Insert(key, value string) error {
	return s.db.Put([]byte(key), []byte(value), nil)
}

//Query查询一条数据.
func (s *StoreEngine) Query(key string) (msgLen int, msg string, err error) {
	b, err := s.db.Get([]byte(key), nil)
	if err != nil && err != leveldb.ErrNotFound {
		return 0, msg, err
	}

	if err == leveldb.ErrNotFound {
		return 0, msg, nil
	}

	msg = string(b)
	return len(msg), msg, nil
}
