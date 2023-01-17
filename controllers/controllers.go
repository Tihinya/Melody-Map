package controllers

import (
	"fmt"
	"groupie-tracker/db"
	"html/template"
	"net/http"
)

// DB

type Card struct {
	Image        string
	GroupName    string
	CreationDate int
	Members      []string
}

type MainData struct {
	Cards []Card
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/html/main-page/index.html", "src/html/main-page/card.html")

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}

	md := MainData{
		Cards: make([]Card, 0),
	}

	for _, artist := range db.DB.Artists {
		card := Card{
			Image:        artist.Image,
			GroupName:    artist.Name,
			CreationDate: artist.CreationDate,
			Members:      artist.Members,
		}

		md.Cards = append(md.Cards, card)
	}

	fmt.Println(md.Cards[0].Image)

	err = t.Execute(w, md)

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}
}

func FullInfo(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/html/full-info/index.html")

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}

	err = t.Execute(w, nil)

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}
}
