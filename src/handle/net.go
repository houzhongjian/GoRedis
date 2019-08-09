package handle

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/houzhongjian/GoRedis/src/conf"
)

//RedisHandle .
type RedisHandle struct {
	Addr string
	Conn net.Conn
	Msg  []string
	Auth bool
}

//NewRedis .
func NewRedis(addr ...string) {
	if len(addr) < 1 {
		addr = append(addr, fmt.Sprintf(":%s", conf.GetString("port")))
	}
	redis := RedisHandle{
		Addr: addr[0],
	}
	redis.Start()
}

func (r *RedisHandle) Start() {
	listen, err := net.Listen("tcp", r.Addr)
	if err != nil {
		log.Printf("%+v\n", err)
		return
	}

	for {
		r.Conn, err = listen.Accept()
		if err != nil {
			log.Printf("%+v\n", err)
			return
		}

		go r.Handle()
	}

}

func (r *RedisHandle) Handle() {
	defer r.Conn.Close()
	for {
		data := make([]byte, 1024)
		_, err := r.Conn.Read(data)
		if err != nil {
			if err == io.EOF {
				r.Conn.Close()
				return
			}
			log.Printf("%+v\n", err)
		}

		r.ParseProtocol(string(data))
	}
}

func (r *RedisHandle) ResponseMsg(msg string) {
	msg = "+" + msg + "\r\n"
	_, err := r.Conn.Write([]byte(msg))
	if err != nil {
		log.Printf("err:%+v\n", err)
		return
	}
}

func (r *RedisHandle) ResponseError(msg string) {
	msg = "-" + msg + "\r\n"
	_, err := r.Conn.Write([]byte(msg))
	if err != nil {
		log.Printf("err:%+v\n", err)
		return
	}
}
