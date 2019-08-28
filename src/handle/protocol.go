package handle

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/houzhongjian/GoRedis/src/utils"

	"github.com/houzhongjian/GoRedis/src/conf"
	"github.com/houzhongjian/GoRedis/src/store"
)

type Protocol struct {
	Conn     net.Conn
	Msg      []string
	Session  map[string]*RedisSession
	Store    *store.StoreEngine
	Stores   []*store.StoreEngine
	IP       string
	SelectDB int //选择的数据库.
}

func (p *Protocol) ResponseMsg(msg interface{}) {
	m := fmt.Sprintf(`"+%v\r\n"`, msg)
	_, err := p.Conn.Write([]byte(m))
	if err != nil {
		log.Printf("err:%+v\n", err)
		return
	}
}

//Select .
func (p *Protocol) Select() {
	if len(p.Msg) != 6 {
		p.Error("(error) ERR wrong number of arguments for 'select' command")
		return
	}

	dbnum, err := utils.ParseInt(p.Msg[4])
	if err != nil {
		p.Error("(error) ERR invalid DB index")
		return
	}
	if dbnum > conf.GetInt("databases") || dbnum < 0 {
		p.Error("(error) ERR DB index is out of range")
		return
	}

	//切换db
	p.SelectDB = dbnum
	p.Store = p.Stores[p.SelectDB]
	p.Success("ok")
}

func (p *Protocol) parse(msg string) {
	proto := strings.Replace(msg, "\r", "", -1)
	p.Msg = strings.Split(proto, "\n")
}

//parseProtocol 解析协议.
func (p *Protocol) ParseProtocol(msg string) {
	p.parse(msg)
	cmd := strings.ToLower(p.Msg[2])
	//连接.
	if p.Connection() {
		p.Success("OK")
		return
	}

	//授权.
	if p.Authorization() {
		p.Auth()
		return
	}

	//检查授权.
	if !p.CheckAuth() {
		p.Error("NOAUTH Authentication required.")
		return
	}

	//检查command是否存在.
	if !p.CheckCommandIsExist(cmd) {
		p.Error("ERR unknown command '" + cmd + "'")
		return
	}

	switch cmd {
	case "set":
		p.Set()
	case "get":
		p.Get()
	case "select":
		p.Select()
	case "info":
		p.Info()
	default:
		log.Println("end")
	}
}

func (p *Protocol) Info() {

}

func (p *Protocol) Auth() {
	if len(p.Msg) != 6 {
		p.Error("ERR wrong number of arguments for 'auth' command")
		return
	}

	pass := p.Msg[4]
	if pass != conf.GetString("requirepass") {
		p.Error("ERR invalid password.")
		return
	}
	ip := p.Conn.RemoteAddr().String()
	p.Session[ip].Auth = true
	p.Success("OK")
}

func (p *Protocol) Set() {
	if len(p.Msg) != 8 {
		p.Error("ERR wrong number of arguments for 'set' command")
		return
	}
	k := p.Msg[4]
	v := p.Msg[6]
	if err := p.Store.Insert(k, v); err != nil {
		p.Error(err.Error())
		return
	}
	p.Success("OK")
}

func (p *Protocol) Connection() bool {
	if p.Msg[2] == "command" || p.Msg[2] == "COMMAND" {
		return true
	}
	return false
}

func (p *Protocol) Get() {
	if len(p.Msg) != 6 {
		p.Error("ERR wrong number of arguments for 'get' command")
		return
	}
	k := p.Msg[4]

	n, msg, err := p.Store.Query(k)
	if err != nil {
		p.Error(err.Error())
		return
	}
	if n < 1 {
		p.Success("(nil)")
		return
	}

	p.ResponseOneMsg(msg)
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
	ip := p.IP
	if conf.IsExist("requirepass") && !p.Session[ip].Auth {
		return false
	}

	return true
}

//CheckCommandIsExist 检查command是否存在.
func (p *Protocol) CheckCommandIsExist(cmd string) bool {
	commandList := []string{
		"auth",
		"set",
		"get",
		"select",
		// "smembers",
	}
	for _, v := range commandList {
		if v == cmd {
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

//SelectDatabase .
func (p *Protocol) SelectDatabase() bool {
	if p.Msg[2] == "select" || p.Msg[2] == "SELECT" {
		return true
	}
	return false
}
