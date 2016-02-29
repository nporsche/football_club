package main

import (
	"errors"
	"html/template"
	"math"
	"net/http"
)

type Revenue struct {
	Time    string
	Amount  int
	Comment string
}

type Match struct {
	Id         int
	Time       string
	Competitor string
	Status     string //到场，缺席，伤病
	Duration   int
	Goals      int
	Assists    int
	Cost       int
	TotalCost  int
}

type PlayerDetail struct {
	Name     string
	Revenues []*Revenue
	Matches  []*Match
}

func accountQueryHandler(w http.ResponseWriter, req *http.Request) {
	player := req.FormValue("name")
	detail, err := getPlayerDetailByName(player)
	if err != nil {
		fillError("get player detail fatal", err, w)
		return
	}
	tmpl, err := template.ParseFiles(*tmplPath + "player.tmpl")
	if err != nil {
		fillError("template parse file fatal", err, w)
		return
	}

	err = tmpl.Execute(w, detail)
	if err != nil {
		fillError("template execute fatal", err, w)
		return
	}
}

func getPlayerDetailByName(player string) (detail *PlayerDetail, err error) {
	detail = &PlayerDetail{
		Name:     player,
		Revenues: make([]*Revenue, 0),
		Matches:  make([]*Match, 0),
	}
	var playerId int
	err = db.QueryRow("select id from players where name=?", player).Scan(&playerId)
	if err != nil {
		return nil, errors.New("无法找到此队员信息" + err.Error())
	}

	rows, err := db.Query("select datetime, amount, reason from revenue_log where player_id = ?", playerId)
	if err != nil {
		return nil, errors.New("查找充值明细异常" + err.Error())
	}

	//收入
	for {
		if !rows.Next() {
			break
		}
		revenue := &Revenue{}
		rows.Scan(&revenue.Time, &revenue.Amount, &revenue.Comment)
		detail.Revenues = append(detail.Revenues, revenue)
	}

	rows, _ = db.Query("select id, datetime, competitor, cost from match_log")
	for {
		if !rows.Next() {
			break
		}
		match := &Match{}
		rows.Scan(&match.Id, &match.Time, &match.Competitor, &match.TotalCost)
		err = fillPlayerMatchInfo(match.Id, playerId, match)
		if err != nil {
			continue
		}

		detail.Matches = append(detail.Matches, match)
	}
	return detail, nil
}

func fillPlayerMatchInfo(matchId, playerId int, match *Match) (err error) {
	var status int
	err = db.QueryRow("select duration, status from duration_log where player_id=? AND match_id=?", playerId, matchId).Scan(&match.Duration, &status)
	match.Status = MatchStatusMap[status]
	if err != nil {
		return err
	} else {
		if status == 0 {
			//正常扣款
			var absence int
			db.QueryRow("select count(*) as absence from duration_log where match_id=? AND status=1", matchId).Scan(&absence)

			var totalDur int
			db.QueryRow("select sum(duration) from duration_log where match_id=? and status=0", matchId).Scan(&totalDur)
			match.Cost = int(math.Ceil((float64(match.TotalCost-10*absence) / float64(totalDur) * float64(match.Duration))))

		} else if status == 1 {
			match.Cost = 10
		} else if status == 2 {
			match.Cost = 0
		}
	}

	//TODO tech

	return nil
}
