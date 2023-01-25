package controllers

import (
	"fmt"
	"groupie-tracker/db"
	"html/template"
	"net/http"
	"sort"
)

// DB

type Card struct {
	FirstAlbum   string
	Location     []string
	Image        string
	GroupName    string
	CreationDate int
	Members      []string
}

type MainData struct {
	Cards        []Card
	CountMembers []int
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/html/main-page/index.html", "src/html/main-page/card.html")

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}

	md := &MainData{
		Cards:        []Card{},
		CountMembers: make([]int, 0),
	}

	temp := make(map[int]int)

	for _, artist := range db.DB.GetArtists() {
		num := len(artist.Members)

		temp[num] = 0

		card := Card{
			FirstAlbum:   artist.FirstAlbum,
			Image:        artist.Image,
			GroupName:    artist.Name,
			CreationDate: artist.CreationDate,
			Members:      artist.Members,
		}

		md.Cards = append(md.Cards, card)
	}

	for k := range temp {
		md.CountMembers = append(md.CountMembers, k)
	}

	sort.Ints(md.CountMembers)

	for i, location := range db.DB.GetLocations() {

		md.Cards[i].Location = location.Location

	}

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

func Search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	query := r.FormValue("query")
	fmt.Println(query)

}
