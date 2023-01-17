package controllers

import (
	"fmt"
	"html/template"
	"net/http"
)

type MainData struct {
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/html/index.html")

	if err != nil {
		fmt.Println(http.StatusInternalServerError)
	}

	err = t.Execute(w, "")

	if err != nil {
		fmt.Println(http.StatusInternalServerError)
	}
}
