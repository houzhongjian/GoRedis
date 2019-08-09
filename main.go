package main

import (
	"log"

	"github.com/houzhongjian/GoRedis/src/command"
	"github.com/houzhongjian/GoRedis/src/handle"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//处理命令行参数.
	command.ParseArgs()

	//启动服务.
	handle.NewRedis()
}
