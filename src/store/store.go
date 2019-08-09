package store

//todo 存储以后采用rocksdb，目前采用map方便开发测试.
var Data = make(map[string]interface{})
