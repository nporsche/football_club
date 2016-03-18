package main

import (
	"errors"
)

var AuthCode = map[string]string{
	"13522759570_hhh": "杭晨",
	"13520511173_x01": "谢永杰",
}

func CheckAuth(code string) (author string, err error) {
	var ok bool
	if author, ok = AuthCode[code]; !ok {
		err = errors.New("Invalid Auth Code")
	}
	return
}
