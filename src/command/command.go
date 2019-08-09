package command

import (
	"os"
)

//version .
func Version() {
	println("GoRedis Version 0.0.1")
	os.Exit(0)
}

//ParseArgs .
func ParseArgs() {
	argv := os.Args
	argc := len(os.Args)
	if argc >= 2 {
		if argv[1] == "-v" || argv[1] == "--version" {
			Version()
		}
	}
}
