package main

import (
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var db *sql.DB
var dbString = flag.String("DB", "ty:ty789@tcp(nporsche.com:3306)/football", "db string")
var matchPath = flag.String("match", "./match_result.yaml", "match result file")

func dbInit() {
	var err error
	db, err = sql.Open("mysql", *dbString)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	flag.Parse()
	dbInit()
	http.HandleFunc("/accountQuery", accountQueryHandler)
	http.HandleFunc("/addMatch", addMatchHandler)
	http.ListenAndServe(":8080", nil)
}
