package handle

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/houzhongjian/GoRedis/src/conf"
	"github.com/houzhongjian/GoRedis/src/store"
)

//RedisHandle .
type RedisHandle struct {
	Lock     sync.RWMutex
	Addr     string
	Store    *store.StoreEngine
	Protocol *Protocol
	Session  int
}

//RedisSession .
type RedisSession struct {
	IP   string
	Auth bool
}

//NewRedis .
func NewRedis(addr ...string) {
	if len(addr) < 1 {
		addr = append(addr, fmt.Sprintf(":%s", conf.GetString("port")))
	}
	redis := RedisHandle{
		Addr:  addr[0],
		Store: store.New(),
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
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("%+v\n", err)
			return
		}

		go r.Handle(conn)
	}
}

func (r *RedisHandle) Handle(conn net.Conn) {
	defer conn.Close()

	ip := conn.RemoteAddr().String()
	proto := Protocol{
		Conn:    conn,
		Store:   r.Store,
		Session: make(map[string]*RedisSession),
	}

	proto.Session[ip] = &RedisSession{
		IP:   ip,
		Auth: false,
	}
	r.Session += 1

	for {
		log.Println("当前连接数:", r.Session)
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				r.Session -= 1
				log.Println(ip, "===>断开连接! 当前连接数为 ", r.Session)
				conn.Close()
				return
			}
			log.Printf("%+v\n", err)
		}

		proto.ParseProtocol(string(buffer))
	}
}
