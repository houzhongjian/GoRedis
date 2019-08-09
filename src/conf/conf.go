package conf

import "log"

func init() {
	if err := load("./conf/redis.conf"); err != nil {
		log.Printf("%+v\n", err)
		return
	}
}
