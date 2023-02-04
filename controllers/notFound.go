package controllers

import (
	"groupie-tracker/errorsSafe"
	"html/template"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/html/404/index.html")

	if err != nil {
		errorsSafe.WrapError(err, errorsSafe.ErrServer)
		return
	}

	err = t.Execute(w, nil)

	if err != nil {
		errorsSafe.WrapError(err, errorsSafe.ErrServer)
		return
	}
}
