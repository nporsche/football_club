package main

import (
	"fmt"
	goyaml "github.com/nporsche/goyaml"
	"io/ioutil"
	"net/http"
	"strings"
)

type MatchInfo struct {
	Datetime   string
	Competitor string
	Cost       int
	Goal       int
	Loss       int
}

type MatchResult struct {
	Match    MatchInfo
	Goal     []string
	Duration []string
}

var result MatchResult

func loadMatchLog(content []byte) error {
	return goyaml.Unmarshal(content, &result)
}

func addMatchHandler(w http.ResponseWriter, req *http.Request) {
	content, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(content))
	if loadMatchLog(content) != nil {
		w.Write([]byte("比赛结果格式错误"))
		return
	}

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		w.Write([]byte("db Begin error"))
		return
	}

	r, err := tx.Exec("INSERT INTO match_log(datetime,competitor,cost,goal,loss) VALUES(?,?,?,?,?)", result.Match.Datetime, result.Match.Competitor, result.Match.Cost, result.Match.Goal, result.Match.Loss)
	if err != nil {
		tx.Rollback()
		w.Write([]byte("insert into match_log error:" + err.Error()))
		return
	}

	matchId, err := r.LastInsertId()
	if err != nil {
		tx.Rollback()
		w.Write([]byte("match_log lastInsertedId error:" + err.Error()))
		return
	}

	for _, goal := range result.Goal {
		g := strings.Split(goal, " ")
		player := g[0]
		goalType := g[1]
		var playerId int
		err := tx.QueryRow("select id from players where name=?", player).Scan(&playerId)
		if err != nil {
			tx.Rollback()
			w.Write([]byte("goal cannot find id of " + player))
			return
		}

		_, err = tx.Exec("INSERT INTO goal_log(match_id,player_id,goal_type) VALUES(?,?,?)", matchId, playerId, goalType)
		if err != nil {
			tx.Rollback()
			w.Write([]byte("INSERT INTO goal_log for " + player + " error:" + err.Error()))
			return
		}
	}

	for _, durationString := range result.Duration {
		durationArray := strings.Split(durationString, " ")
		player := durationArray[0]
		dur := durationArray[1]
		status := durationArray[2]
		var playerId int
		err := tx.QueryRow("select id from players where name=?", player).Scan(&playerId)
		if err != nil {
			tx.Rollback()
			w.Write([]byte("duration cannot find id of " + player))
			return
		}

		_, err = tx.Exec("INSERT INTO duration_log(match_id,player_id,duration, status) VALUES(?,?,?,?)", matchId, playerId, dur, status)
		if err != nil {
			tx.Rollback()
			w.Write([]byte("INSERT INTO duration_log for " + player + " error:" + err.Error()))
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		w.Write([]byte("Commit error:" + err.Error()))
		return
	}

	w.Write([]byte("MATCH ADDED!"))
}
