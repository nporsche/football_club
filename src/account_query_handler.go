package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"math"
	"net/http"
)

var statusStringMap map[int]string

func init() {
	statusStringMap = make(map[int]string)
	statusStringMap[0] = "正常"
	statusStringMap[1] = "伤病"
	statusStringMap[2] = "退出"
}

func accountQueryHandler(w http.ResponseWriter, req *http.Request) {
	player := req.FormValue("name")
	var playerId int
	var status int
	err := db.QueryRow("select id,status from players where name=?", player).Scan(&playerId, &status)
	if err != nil {
		w.Write([]byte("该队员不存在,请联系队长谢永杰"))
		return
	}
	sum := 0
	err = db.QueryRow("select sum(amount) as sum from revenue_log where player_id=?", playerId).Scan(&sum)
	if err != nil {
		w.Write([]byte("该队员无充值记录"))
		return
	}

	disp := bytes.NewBufferString(fmt.Sprintf("姓名：%s\n", player))
	disp.WriteString(fmt.Sprintf("状态：%s\n\n", statusStringMap[status]))
	rows, err := db.Query("select datetime, amount, reason from revenue_log where player_id = ?", playerId)
	if err != nil {
		w.Write([]byte("查询充值明细异常"))
		return
	}

	//收入
	disp.WriteString("收入表\n")
	disp.WriteString("--------------------------------------------------------------\n")
	disp.WriteString("充值时间			充值金额		备注\n")
	for {
		if !rows.Next() {
			break
		}
		var date string
		var amount int
		var reason string
		rows.Scan(&date, &amount, &reason)
		disp.WriteString(fmt.Sprintf("%s	%d		%s\n", date, amount, reason))
	}

	disp.WriteString("--------------------------------------------------------------\n")
	disp.WriteString(fmt.Sprintf("充值总金额: %d 元\n", sum))
	//
	disp.WriteString("\n")
	//支出
	disp.WriteString("支出表\n")
	disp.WriteString("---------------------------------------------------------------------------------------------------\n")
	disp.WriteString("比赛时间			对手		上场时间		备注		消费金额\n")
	rows, _ = db.Query("select id, datetime, competitor, cost from match_log")
	totalCost := 0
	for {
		var disDatetime string
		var disCompetitor string
		var disDuration int
		var disStatus string
		var disCost int

		if !rows.Next() {
			break
		}
		var matchId int
		var cost int
		rows.Scan(&matchId, &disDatetime, &disCompetitor, &cost)

		var status int
		err = db.QueryRow("select duration, status from duration_log where player_id=? AND match_id=?", playerId, matchId).Scan(&disDuration, &status)
		if err == nil {
			if status == 0 {
				//正常扣款
				var absence int
				db.QueryRow("select count(*) as absence from duration_log where match_id=? AND status=1", matchId).Scan(&absence)

				var totalDur int
				db.QueryRow("select sum(duration) from duration_log where match_id=? and status=0", matchId).Scan(&totalDur)
				disCost = int(math.Ceil((float64(cost-10*absence) / float64(totalDur) * float64(disDuration))))
				disStatus = "正常"

			} else if status == 1 {
				disStatus = "缺勤"
				disCost = 10
			} else if status == 2 {
				disStatus = "伤病"
				disCost = 0
			}
			totalCost += disCost
			disp.WriteString(fmt.Sprintf("%s	%s		%d分钟		%s		%d\n", disDatetime, disCompetitor, disDuration, disStatus, disCost))
		} else if err == sql.ErrNoRows {
			continue
		} else {
			disp.WriteString("错误，请刷新!\n")
			break
		}
	}
	disp.WriteString("---------------------------------------------------------------------------------------------------\n")
	disp.WriteString(fmt.Sprintf("总支出：%d 元\n\n", totalCost))
	disp.WriteString(fmt.Sprintf("================\n账户余额：%d 元\n================\n", sum-totalCost))
	//
	disp.WriteString("\n\n\n")
	//技术统计
	disp.WriteString("技术统计\n")
	disp.WriteString("-----------------------\n")
	disp.WriteString("进球类型		数量\n")
	rows, err = db.Query("select goal_type,count(goal_type) from goal_log where player_id=? group by goal_type", playerId)
	if err == nil {
		for {
			if !rows.Next() {
				break
			}
			var goalType string
			var count int
			rows.Scan(&goalType, &count)
			disp.WriteString(fmt.Sprintf("%s		%d\n", goalType, count))
		}
	}
	disp.WriteString("-----------------------\n")

	w.Write(disp.Bytes())
}
