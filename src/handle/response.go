package handle

import (
	"fmt"
	"log"
)

//Success 返回一个成功的状态信息.
func (p *Protocol) Success(msg string) {
	m := fmt.Sprintf("+%s\r\n", msg)
	_, err := p.Conn.Write([]byte(m))
	if err != nil {
		log.Printf("err:%+v\n", err)
		return
	}
}

//ResponseError 返回一个错误信息.
func (p *Protocol) Error(msg string) {
	m := fmt.Sprintf("-%s\r\n", msg)
	_, err := p.Conn.Write([]byte(m))
	if err != nil {
		log.Printf("err:%+v\n", err)
		return
	}
}

func (p *Protocol) ResponseOneMsg(msg string) {
	m := fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg)
	_, err := p.Conn.Write([]byte(m))
	if err != nil {
		log.Printf("err:%+v\n", err)
		return
	}
}
