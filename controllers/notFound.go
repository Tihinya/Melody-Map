package controllers

import (
	"fmt"
	"html/template"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/html/404/index.html")

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}

	err = t.Execute(w, nil)

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}
}
