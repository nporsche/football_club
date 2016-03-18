package main

import (
	"html/template"
	"net/http"
)

func addMatchPage(w http.ResponseWriter, req *http.Request) {
	type Player struct {
		Name string
		Tag  int
	}
	rows, err := db.Query("select name, tag from players where status=0 OR status=1")
	if err != nil {
		fillError("cannot select players", err, w)
		return
	}
	players := []*Player{}
	for {
		if !rows.Next() {
			break
		}
		var name string
		var tag int
		rows.Scan(&name, &tag)
		player := &Player{Name: name,
			Tag: tag}
		players = append(players, player)
	}
	tmpl, err := template.ParseFiles(*tmplPath + "addMatchPage.tmpl")
	if err != nil {
		fillError("template parse file fatal", err, w)
		return
	}

	err = tmpl.Execute(w, players)
	if err != nil {
		fillError("template execute fatal", err, w)
		return
	}
}
