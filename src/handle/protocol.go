package handle

import (
	"log"
	"strings"

	"github.com/houzhongjian/GoRedis/src/store"
)

//parseProtocol 解析协议.
func (r *RedisHandle) ParseProtocol(msg string) {
	proto := strings.Replace(msg, "\r", "", -1)
	r.Msg = strings.Split(proto, "\n")
	log.Println(r.Msg)
	if r.Login() {
		r.ResponseMsg("OK")
	}

	if r.Set() {
		k := r.Msg[4]
		v := r.Msg[6]
		store.Data[k] = v
		r.ResponseMsg("OK")
	}

	if r.Get() {
		k := r.Msg[4]
		v := store.Data[k]
		r.ResponseMsg(v)
	}
}

func (r *RedisHandle) Login() bool {
	if r.Msg[2] == "command" || r.Msg[2] == "COMMAND" {
		return true
	}
	return false
}

func (r *RedisHandle) Set() bool {
	if r.Msg[2] == "set" || r.Msg[2] == "SET" {
		return true
	}
	return false
}

func (r *RedisHandle) Get() bool {
	if r.Msg[2] == "get" || r.Msg[2] == "GET" {
		return true
	}
	return false
}
