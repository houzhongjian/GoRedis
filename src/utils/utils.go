package utils

import "strconv"

func ParseInt(s string) (num int, err error) {
	num, err = strconv.Atoi(s)
	return num, err
}
