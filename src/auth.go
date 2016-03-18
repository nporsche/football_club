package main

import (
	"errors"
)

var AuthCode = map[string]string{
	"hc_1250":  "杭晨",
	"xyj_1251": "谢永杰",
	"lb_1252":  "吕博",
	"cpy_1253": "陈鹏宇",
	"czs_1254": "陈振升",
	"wyx_1255": "王英鑫",
}

func CheckAuth(code string) (author string, err error) {
	var ok bool
	if author, ok = AuthCode[code]; !ok {
		err = errors.New("Invalid Auth Code")
	}
	return
}
