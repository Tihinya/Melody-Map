package controllers

import (
	"groupie-tracker/db"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Card struct {
	Id           int
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
	GroupNames   []string
	CreationDate []int
	FirstAlbum   []string
	Members      []string
	Locations    []string
}

// TODO:
// 0. logo --> home page 				+
// 1. fix fonts
// 2. show creation date in header
// 3. error handling with Rick Ashley
// 4. 404 ^^^^^^^^^^^^^^^^^^^^^^^^^^
// 5. show first album in full page

func prepareFilters(md *MainData, temp map[int]interface{}) {
	locations := db.DB.GetLocations()

	uniqueLocations := make(map[string]interface{})
	uniqueMembers := make(map[string]interface{})
	uniqueDate := make(map[int]interface{})
	uniqueFirstAlbum := make(map[string]interface{})

	for _, v := range locations {
		for _, v1 := range v.Location {
			uniqueLocations[v1] = nil
		}
	}

	artist := db.DB.GetArtists()
	for _, v := range artist {
		for _, v1 := range v.Members {
			uniqueMembers[v1] = nil
		}
		uniqueDate[v.CreationDate] = nil
		md.GroupNames = append(md.GroupNames, v.Name)
		uniqueFirstAlbum[v.FirstAlbum] = nil
	}

	for k := range temp {
		md.CountMembers = append(md.CountMembers, k)
	}

	md.CountMembers = mapToSlice(temp)
	md.Members = mapToSlice(uniqueMembers)
	md.Locations = mapToSlice(uniqueLocations)
	md.CreationDate = mapToSlice(uniqueDate)
	md.FirstAlbum = mapToSlice(uniqueFirstAlbum)

	sort.Ints(md.CountMembers)
	sort.Strings(md.Locations)
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/html/main-page/index.html", "src/html/main-page/card.html")

	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. We are working on that", http.StatusInternalServerError)
		return
	}

	md := &MainData{
		Cards:        []Card{},
		CountMembers: make([]int, 0),
		Locations:    make([]string, 0),
	}

	temp := make(map[int]interface{})
	search := strings.ToLower(r.FormValue("search-input"))

	filter := getFilters(r)

	locations := db.DB.GetLocations()
	for i, artist := range db.DB.GetArtists() {
		num := len(artist.Members)
		temp[num] = nil

		a, _ := time.Parse("02-01-2006", artist.FirstAlbum)

		year := a.Year()

		if (filter.CreationDateStart != 0 || filter.CreationDateEnd != 0) &&
			(artist.CreationDate < filter.CreationDateStart || artist.CreationDate > filter.CreationDateEnd) {
			continue
		}

		if (filter.AlbumDateStart != 0 || filter.AlbumDateEnd != 0) &&
			(year < filter.AlbumDateStart || year > filter.AlbumDateEnd) {
			continue
		}

		if len(filter.Members) > 0 && !contains(filter.Members, len(artist.Members)) {
			continue
		}

		if filter.Location != "" && !contains(locations[i].Location, filter.Location) {
			continue
		}

		isArtist(artist, search, locations[i].Location)

		card := Card{
			Id:           artist.Id,
			FirstAlbum:   artist.FirstAlbum,
			Image:        artist.Image,
			GroupName:    artist.Name,
			Location:     locations[artist.Id-1].Location,
			CreationDate: artist.CreationDate,
			Members:      artist.Members,
		}

		md.Cards = append(md.Cards, card)
	}

	prepareFilters(md, temp)

	err = t.Execute(w, md)

	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong. We are working on that", http.StatusInternalServerError)
		return
	}
}
