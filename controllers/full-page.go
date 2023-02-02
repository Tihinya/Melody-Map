package controllers

import (
	"fmt"
	"groupie-tracker/db"
	"groupie-tracker/router"
	"html/template"
	"net/http"
	"strconv"
)

func FullInfo(w http.ResponseWriter, r *http.Request) {
	sid := router.GetField(r, "id")

	id, err := strconv.Atoi(sid)
	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}

	md := Info{
		Card:          Card{},
		LocationDates: map[string][]string{},
	}

	for _, artist := range db.DB.GetArtists() {
		if artist.Id == id {
			card := Card{
				Id:           artist.Id,
				Image:        artist.Image,
				GroupName:    artist.Name,
				CreationDate: artist.CreationDate,
				Members:      artist.Members,
			}

			md.Card = card
			break
		}
	}

	for _, dl := range db.DB.GetRelations() {
		if dl.Id == id {
			for k, v := range dl.DatesLocations {
				md.LocationDates[k] = v
			}
		}
	}

	t, err := template.ParseFiles("src/html/full-info/index.html")

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}

	err = t.Execute(w, md)

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}
}
