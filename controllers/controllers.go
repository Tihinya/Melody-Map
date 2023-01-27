package controllers

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/db"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Card struct {
	Image        string
	GroupName    string
	CreationDate int
	Members      []string
}

type Info struct {
	Card          Card
	LocationDates map[string][]string
}

type MainData struct {
	Cards []Card
}

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/html/main-page/index.html", "src/html/main-page/card.html")

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}

	md := MainData{
		Cards: []Card{},
	}

	for _, artist := range db.DB.GetArtists() {
		card := Card{
			Image:        artist.Image,
			GroupName:    artist.Name,
			CreationDate: artist.CreationDate,
			Members:      artist.Members,
		}

		md.Cards = append(md.Cards, card)
	}

	err = t.Execute(w, md)

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}
}

func FullInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	id := 1

	md := Info{
		Card:          Card{},
		LocationDates: map[string][]string{},
	}

	for i, artist := range db.DB.GetArtists() {
		if i == id {
			card := Card{
				Image:        artist.Image,
				GroupName:    artist.Name,
				CreationDate: artist.CreationDate,
				Members:      artist.Members,
			}

			md.Card = card
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

func DatesLocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	k := os.Getenv("GOOGLE_API_KEY")

	id := 1

	var coordinates []Coordinate

	if k == "" {
		return
	}

	for _, dl := range db.DB.GetRelations() {
		if dl.Id != id {
			continue
		}

		for loc := range dl.DatesLocations {
			addr := strings.Replace(loc, "_", "+", -1)

			req := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%v&key=%v", addr, k)
			res, _ := http.Get(req)

			lat, lng := db.GetGoogleMap(res)

			ll := Coordinate{
				Lat: lat,
				Lng: lng,
			}

			coordinates = append(coordinates, ll)
		}

		break
	}

	json.NewEncoder(w).Encode(coordinates)
}
