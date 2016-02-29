package main

import (
	"fmt"
	"net/http"
)

var playerStatusMap = map[int]string{0: "正常", 1: "伤病", 2: "退出"}
var MatchStatusMap = map[int]string{0: "正常", 1: "缺勤", 2: "伤病"}

func fillError(detail string, err error, rw http.ResponseWriter) {
	res := fmt.Sprintf("error:%s\ndetail:%s\n", err.Error(), detail)
	rw.Write([]byte(res))
}
