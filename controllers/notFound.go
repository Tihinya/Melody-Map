package controllers

import (
	"html/template"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/html/404/indexs.html")

	if err != nil {
		http.Error(w, "Something went wrong. We are working on that", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)

	if err != nil {
		http.Error(w, "Something went wrong. We are working on that", http.StatusInternalServerError)
		return
	}
}
