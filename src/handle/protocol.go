package handle

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/houzhongjian/GoRedis/src/conf"
	"github.com/houzhongjian/GoRedis/src/store"
)

type Protocol struct {
	Conn    net.Conn
	Msg     []string
	Session map[string]*RedisSession
	Store   *store.StoreEngine
}

func (p *Protocol) ResponseMsg(msg interface{}) {
	m := fmt.Sprintf("+%v\r\n", msg)
	_, err := p.Conn.Write([]byte(m))
	if err != nil {
		log.Printf("err:%+v\n", err)
		return
	}
}

func (p *Protocol) ResponseError(msg string) {
	msg = "-" + msg + "\r\n"
	_, err := p.Conn.Write([]byte(msg))
	if err != nil {
		log.Printf("err:%+v\n", err)
		return
	}
}

//parseProtocol 解析协议.
func (p *Protocol) ParseProtocol(msg string) {
	proto := strings.Replace(msg, "\r", "", -1)
	p.Msg = strings.Split(proto, "\n")

	//连接.
	if p.Login() {
		p.ResponseMsg("OK")
		return
	}

	//检查command是否存在.
	if !p.CheckCommandIsExist() {
		p.ResponseError("ERR unknown command '" + p.Msg[2] + "'")
		return
	}

	//授权.
	if p.Authorization() {
		if len(p.Msg) != 6 {
			p.ResponseError("ERR wrong number of arguments for 'auth' command")
			return
		}

		pass := p.Msg[4]
		if pass != conf.GetString("requirepass") {
			p.ResponseError("ERR invalid password.")
			return
		}
		ip := p.Conn.RemoteAddr().String()
		p.Session[ip].Auth = true
		p.ResponseMsg("OK")
		return
	}

	//检查授权.
	if !p.CheckAuth() {
		p.ResponseError("NOAUTH Authentication required.")
		return
	}

	if p.Set() {
		if len(p.Msg) != 8 {
			p.ResponseError("ERR wrong number of arguments for 'set' command")
			return
		}
		k := p.Msg[4]
		v := p.Msg[6]
		if err := p.Store.Insert(k, v); err != nil {
			p.ResponseError(err.Error())
			return
		}
		p.ResponseMsg("OK")
		return
	}

	if p.Get() {
		if len(p.Msg) != 6 {
			p.ResponseError("ERR wrong number of arguments for 'get' command")
			return
		}
		k := p.Msg[4]

		// var msg string
		// if _, ok := store.Data[k]; ok {
		// 	msg = store.Data[k].(string)
		// 	vType := reflect.TypeOf(msg).String()
		// 	if vType != "string" {
		// 		p.ResponseError("WRONGTYPE Operation against a key holding the wrong kind of value")
		// 		return
		// 	}
		// }

		b, err := p.Store.Query(k)
		if err != nil {
			p.ResponseError(err.Error())
			return
		}

		p.ResponseMsg(string(b))
		return
	}

	// if p.Sadd() {
	// 	if len(p.Msg) < 8 {
	// 		p.ResponseError("ERR wrong number of arguments for 'sadd' command")
	// 		return
	// 	}
	// 	k := p.Msg[4]

	// 	var v []interface{}
	// 	for i := 5; i < len(p.Msg); i++ {
	// 		if i%2 == 0 {
	// 			log.Println(p.Msg[i])
	// 			v = append(v, p.Msg[i])
	// 		}
	// 	}

	// 	if _, ok := store.Data[k]; ok {
	// 		val := store.Data[k]
	// 		vType := reflect.TypeOf(val).String()
	// 		if vType == "[]interface {}" {
	// 			vv := val.([]interface{})
	// 			for _, item := range vv {
	// 				v = append(v, item)
	// 			}
	// 		}
	// 	}

	// 	store.Data[k] = v
	// 	p.ResponseMsg("OK")
	// }

	// if p.Smembers() {
	// 	log.Println(len(p.Msg))
	// 	if len(p.Msg) != 6 {
	// 		p.ResponseError("ERR wrong number of arguments for 'smembers' command")
	// 		return
	// 	}

	// 	k := p.Msg[4]

	// 	var msg string
	// 	//判断k是否存在.
	// 	if _, ok := store.Data[k]; ok {
	// 		if reflect.TypeOf(store.Data[k]).String() != "[]interface {}" {
	// 			p.ResponseError("WRONGTYPE Operation against a key holding the wrong kind of value")
	// 			return
	// 		}
	// 		v := store.Data[k].([]interface{})
	// 		for k, item := range v {
	// 			msg += fmt.Sprintf("%d) \"%s\"\n", k+1, item)
	// 		}
	// 	}

	// 	p.ResponseMsg(msg)
	// 	return
	// }

	log.Println("end...")
}

func (p *Protocol) Login() bool {
	if p.Msg[2] == "command" || p.Msg[2] == "COMMAND" {
		return true
	}
	return false
}

func (p *Protocol) Set() bool {
	if p.Msg[2] == "set" || p.Msg[2] == "SET" {
		return true
	}
	return false
}

func (p *Protocol) Get() bool {
	if p.Msg[2] == "get" || p.Msg[2] == "GET" {
		return true
	}
	return false
}

//Authorization 登录授权.
func (p *Protocol) Authorization() bool {
	if p.Msg[2] == "auth" || p.Msg[2] == "AUTH" {
		return true
	}
	return false
}

//CheckAuth 验证是否需要进行权限认证.
func (p *Protocol) CheckAuth() bool {
	//判断是否设置了密码.
	ip := p.Conn.RemoteAddr().String()
	if conf.IsExist("requirepass") && !p.Session[ip].Auth {
		return false
	}

	return true
}

//CheckCommandIsExist 检查command是否存在.
func (p *Protocol) CheckCommandIsExist() bool {
	commandList := []string{
		"auth",
		"set",
		"get",
		// "sadd",
		// "smembers",
	}
	for _, v := range commandList {
		if v == strings.ToLower(p.Msg[2]) {
			return true
		}
	}

	return false
}

//Sadd .
func (p *Protocol) Sadd() bool {
	if p.Msg[2] == "sadd" || p.Msg[2] == "SADD" {
		return true
	}
	return false
}

//Smembers .
func (p *Protocol) Smembers() bool {
	if p.Msg[2] == "smembers" || p.Msg[2] == "SMEMBERS" {
		return true
	}
	return false
}
