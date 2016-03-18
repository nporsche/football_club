package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type MatchInfoV2 struct {
	Datetime   string `json:"datetime"`
	Competitor string `json:"competitor"`
	Amount     int    `json:"amount"`
	Goal       int    `json:"goal"`
	Loss       int    `json:"loss"`
}

type MatchPlayerInfoV2 struct {
	Name     string `json:"name"`
	Tag      int    `json:"tag"`
	Status   int    `json:"status"`
	Duration int    `json:"duration"`
	FreeKick int    `json:"freekick"` //定位球
	Penalty  int    `json:"penalty"`  //点球
	Mobile   int    `json:"mobile"`   //运动战
	Assist   int    `json:"assist"`   //助攻
}

type MatchResultV2 struct {
	Match        MatchInfoV2          `json:"match"`
	MatchPlayers []*MatchPlayerInfoV2 `json:"playerInfo"`
	AuthCode     string               `json:"authCode"`
}

func loadMatchLogV2(content []byte) (result MatchResultV2, err error) {
	err = json.Unmarshal(content, &result)
	return
}

func addMatchHandlerV2(w http.ResponseWriter, req *http.Request) {
	content, _ := ioutil.ReadAll(req.Body)
	result, err := loadMatchLogV2(content)
	if err != nil {
		w.Write([]byte("比赛结果格式错误" + err.Error()))
		return
	}
	author, err := CheckAuth(result.AuthCode)
	if err != nil {
		w.Write([]byte("授权码有误,非法操作"))
		return
	}

	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		w.Write([]byte("db Begin error"))
		return
	}

	r, err := tx.Exec("INSERT INTO match_log(datetime,competitor,cost,goal,loss,author) VALUES(?,?,?,?,?,?)", result.Match.Datetime, result.Match.Competitor, result.Match.Amount, result.Match.Goal, result.Match.Loss, author)
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

	for _, playerInfo := range result.MatchPlayers {
		var playerId int
		err := tx.QueryRow("select id from players where name=? AND tag=?", playerInfo.Name, playerInfo.Tag).Scan(&playerId)
		if err != nil {
			tx.Rollback()
			w.Write([]byte("duration cannot find id of " + playerInfo.Name))
			return
		}

		_, err = tx.Exec("INSERT INTO duration_log(match_id,player_id,duration, status, author) VALUES(?,?,?,?,?)", matchId, playerId, playerInfo.Duration, playerInfo.Status, author)
		if err != nil {
			tx.Rollback()
			w.Write([]byte("INSERT INTO duration_log for " + playerInfo.Name + " error:" + err.Error()))
			return
		}

		err = insertTechData(tx, int(matchId), playerId, playerInfo.FreeKick, playerInfo.Penalty, playerInfo.Mobile, playerInfo.Assist, author)
		if err != nil {
			tx.Rollback()
			w.Write([]byte("INSERT INTO goal_log for " + playerInfo.Name + " error:" + err.Error()))
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		w.Write([]byte("Commit error:" + err.Error()))
		return
	}

	w.Write([]byte("添加比赛成功，作者:" + author))
}

func insertTechData(tx *sql.Tx, matchId, playerId, freeKick, penalty, mobile, assist int, author string) (err error) {
	for i := 0; i < freeKick; i++ {
		_, err = tx.Exec("INSERT INTO goal_log(match_id,player_id,goal_type,author) VALUES(?,?,?,?)", matchId, playerId, "任意球", author)
		if err != nil {
			return err
		}
	}

	for i := 0; i < penalty; i++ {
		_, err = tx.Exec("INSERT INTO goal_log(match_id,player_id,goal_type,author) VALUES(?,?,?,?)", matchId, playerId, "点球", author)
		if err != nil {
			return err
		}
	}

	for i := 0; i < mobile; i++ {
		_, err = tx.Exec("INSERT INTO goal_log(match_id,player_id,goal_type,author) VALUES(?,?,?,?)", matchId, playerId, "运动战", author)
		if err != nil {
			return err
		}
	}
	for i := 0; i < assist; i++ {
		_, err = tx.Exec("INSERT INTO goal_log(match_id,player_id,goal_type,author) VALUES(?,?,?,?)", matchId, playerId, "助攻", author)
		if err != nil {
			return err
		}
	}
	return nil
}
