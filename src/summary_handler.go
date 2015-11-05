package main

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
)

func summaryHandler(w http.ResponseWriter, req *http.Request) {
	disp := bytes.NewBufferString("----------------------\n")
	disp.WriteString("姓名		余额\n")
	rows, err := db.Query("select id, name from players")
	if err != nil {
		w.Write([]byte("查询球员名单错误"))
		return
	}
	for {
		if !rows.Next() {
			break
		}
		var name string
		var playerId int
		rows.Scan(&playerId, &name)
		b, e := getAccountByPlayerId(playerId)
		if e == nil {
			disp.WriteString(fmt.Sprintf("%s		%d\n", name, b))
		} else {
			disp.WriteString(fmt.Sprintf("%s		%s\n", name, "异常"))
		}
	}
	disp.WriteString("----------------------\n")
	w.Write(disp.Bytes())
}
func getAccountByPlayerId(playerId int) (balance int, err error) {
	sum := 0
	err = db.QueryRow("select sum(amount) as sum from revenue_log where player_id=?", playerId).Scan(&sum)
	if err != nil {
		return
	}

	rows, _ := db.Query("select id, datetime, competitor, cost from match_log")
	totalCost := 0
	for {
		var disDatetime string
		var disCompetitor string
		var disDuration int
		var disCost int

		if !rows.Next() {
			break
		}
		var matchId int
		var cost int
		rows.Scan(&matchId, &disDatetime, &disCompetitor, &cost)

		var status int
		db.QueryRow("select duration, status from duration_log where player_id=? AND match_id=?", playerId, matchId).Scan(&disDuration, &status)
		if status == 0 {
			//正常扣款
			var absence int
			db.QueryRow("select count(*) as absence from duration_log where match_id=? AND status=1", matchId).Scan(&absence)

			var totalDur int
			db.QueryRow("select sum(duration) from duration_log where match_id=? and status=0", matchId).Scan(&totalDur)
			disCost = int(math.Ceil((float64(cost-10*absence) / float64(totalDur) * float64(disDuration))))

		} else if status == 1 {
			disCost = 10
		} else if status == 2 {
			disCost = 0
		}
		totalCost += disCost
	}
	balance = sum - totalCost
	return
}
