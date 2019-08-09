package handle

import (
	"log"
	"strings"

	"github.com/houzhongjian/GoRedis/src/conf"
	"github.com/houzhongjian/GoRedis/src/store"
)

//parseProtocol 解析协议.
func (r *RedisHandle) ParseProtocol(msg string) {
	proto := strings.Replace(msg, "\r", "", -1)
	r.Msg = strings.Split(proto, "\n")
	log.Println(r.Msg)

	//连接.
	if r.Login() {
		r.ResponseMsg("OK")
		return
	}

	//检查command是否存在.
	if !r.CheckCommandIsExist() {
		r.ResponseError("ERR unknown command '" + r.Msg[2] + "'")
		return
	}

	//授权.
	if r.Authorization() {
		if len(r.Msg) != 6 {
			r.ResponseError("ERR wrong number of arguments for 'auth' command")
			return
		}

		pass := r.Msg[4]
		if pass != conf.GetString("requirepass") {
			log.Println(conf.GetString("requirepass"))
			r.ResponseError("ERR invalid password.")
			return
		}
		r.Auth = true
		r.ResponseMsg("OK")
	}

	//检查授权.
	if !r.CheckAuth() {
		r.ResponseError("NOAUTH Authentication required.")
		return
	}

	if r.Set() {
		if len(r.Msg) != 8 {
			r.ResponseError("ERR wrong number of arguments for 'set' command")
			return
		}
		k := r.Msg[4]
		v := r.Msg[6]
		store.Data[k] = v
		r.ResponseMsg("OK")
	}

	if r.Get() {
		if len(r.Msg) != 6 {
			r.ResponseError("ERR wrong number of arguments for 'get' command")
			return
		}
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

//Authorization 登录授权.
func (r *RedisHandle) Authorization() bool {
	if r.Msg[2] == "auth" || r.Msg[2] == "AUTH" {
		return true
	}
	return false
}

//CheckAuth 验证是否需要进行权限认证.
func (r *RedisHandle) CheckAuth() bool {
	//判断是否设置了密码.
	if conf.IsExist("requirepass") && !r.Auth {
		return false
	}

	return true
}

//CheckCommandIsExist 检查command是否存在.
func (r *RedisHandle) CheckCommandIsExist() bool {
	commandList := []string{
		"auth",
		"set",
		"get",
	}
	for _, v := range commandList {
		if v == r.Msg[2] {
			return true
		}
	}

	return false
}
