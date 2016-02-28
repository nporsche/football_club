package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var db *sql.DB
var dbString = flag.String("DB", "ty:ty789@tcp(nporsche.com:3306)/football", "db string")
var port = flag.Int("port", 8080, "http port")
var tmplPath = flag.String("template-path", "./static/", "template path")

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
	http.HandleFunc("/summary", summaryHandler)
	http.HandleFunc("/index.html", summaryHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
