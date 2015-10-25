package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	goyaml "github.com/nporsche/goyaml"
	"io/ioutil"
	"log"
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

func loadMatchLog(path string) {
	var content []byte
	var err error
	if content, err = ioutil.ReadFile(path); err != nil {
		log.Fatalln("ReloadConfig failure from path:" + path)
		return
	}
	if err := goyaml.Unmarshal(content, &result); err != nil {
		log.Fatalln("ReloadConfig unmarshal failure from path:" + path)
		return
	}
}

func main() {
	matchPath := flag.String("match", "./match_result.yaml", "match result file")
	loadMatchLog(*matchPath)
	dbString := flag.String("DB", "ty:ty789@tcp(nporsche.com:3306)/football", "db string")
	db, err := sql.Open("mysql", *dbString)
	if err != nil {
		log.Fatalln(err)
	}
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	r, err := tx.Exec("INSERT INTO match_log(datetime,competitor,cost,goal,loss) VALUES(?,?,?,?,?)", result.Match.Datetime, result.Match.Competitor, result.Match.Cost, result.Match.Goal, result.Match.Loss)
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	matchId, err := r.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatalln(err)
	}

	for _, goal := range result.Goal {
		g := strings.Split(goal, " ")
		player := g[0]
		goalType := g[1]
		var playerId int
		err := tx.QueryRow("select id from players where name=?", player).Scan(&playerId)
		if err != nil {
			tx.Rollback()
			log.Fatalln(err)
		}

		_, err = tx.Exec("INSERT INTO goal_log(match_id,player_id,goal_type) VALUES(?,?,?)", matchId, playerId, goalType)
		if err != nil {
			tx.Rollback()
			log.Fatalln(err)
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
			log.Fatalln(err)
		}

		_, err = tx.Exec("INSERT INTO duration_log(match_id,player_id,duration, status) VALUES(?,?,?,?)", matchId, playerId, dur, status)
		if err != nil {
			tx.Rollback()
			log.Fatalln(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("match result added")
}
