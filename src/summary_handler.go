package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
)

const (
	ACTIVE  = "active"
	SUCCESS = "success"
	WARNING = "warning"
	DANGER  = "danger"
)

type TeamSummary struct {
	Revenue int
	Cost    int
	Balance int
}

type PlayerSummary struct {
	Name       string
	Tag        int
	Attendance string
	Balance    int
	Status     string
	ColorState string
}

type Summary struct {
	Team    TeamSummary
	Players []*PlayerSummary
}

func summaryHandler(w http.ResponseWriter, req *http.Request) {
	sum, err := summaryProcess()
	if err != nil {
		fillError("summary process fatal", err, w)
		return
	}
	tmpl, err := template.ParseFiles(*tmplPath + "summary.tmpl")
	if err != nil {
		fillError("template parse file fatal", err, w)
		return
	}

	err = tmpl.Execute(w, sum)
	if err != nil {
		fillError("template execute fatal", err, w)
		return
	}
}

func summaryProcess() (sum *Summary, err error) {
	sum = &Summary{Team: TeamSummary{},
		Players: make([]*PlayerSummary, 0)}

	revenue, err := getTotalRevenue()
	if err != nil {
		return nil, err
	}
	cost, err := getTotalCost()
	if err != nil {
		return nil, err
	}
	sum.Team.Revenue = revenue
	sum.Team.Cost = cost
	sum.Team.Balance = revenue - cost

	rows, err := db.Query("select id, name, tag, status from players")
	if err != nil {
		return nil, err
	}
	for {
		if !rows.Next() {
			break
		}
		var name string
		var playerId int
		var tag int
		var status int
		rows.Scan(&playerId, &name, &tag, &status)
		if status == 2 {
			//已经离队不予处理
			continue
		}
		player := &PlayerSummary{Name: name,
			Tag:        tag,
			Status:     playerStatusMap[status],
			ColorState: SUCCESS,
			Attendance: "异常",
			Balance:    0}

		b, e := getAccountByPlayerId(playerId)
		if e == nil {
			player.Balance = b
		}

		attendance, e := getAttendanceByPlayerId(playerId)
		if e == nil {
			player.Attendance = fmt.Sprintf("%d%%", attendance)
		}

		//color state
		if b < 0 {
			player.ColorState = DANGER
		} else if b < 50 {
			player.ColorState = WARNING
		}
		if status == 1 {
			//伤病
			player.ColorState = ACTIVE
		}
		//color state end

		sum.Players = append(sum.Players, player)
	}
	return sum, nil
}

func getTotalRevenue() (revenue int, err error) {
	err = db.QueryRow("select sum(amount) from revenue_log").Scan(&revenue)
	return
}

func getTotalCost() (cost int, err error) {
	err = db.QueryRow("select sum(cost) from match_log").Scan(&cost)
	return
}

func getAttendanceByPlayerId(playerId int) (attendRate int, err error) {
	var matchCnt int
	err = db.QueryRow("select count(*) from duration_log where player_id=?", playerId).Scan(&matchCnt)
	if err != nil {
		return 0, err
	}
	if matchCnt == 0 {
		return 0, nil
	}
	var available int
	err = db.QueryRow("select count(*) from duration_log where player_id=? AND status=0", playerId).Scan(&available)
	if err != nil {
		return 0, err
	}
	return int((float64(available) / float64(matchCnt)) * 100), nil
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
