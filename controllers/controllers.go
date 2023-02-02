package controllers

import (
	"fmt"
	"groupie-tracker/db"
	"html/template"
	"net/http"
	"sort"
	"strconv"
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

type Info struct {
	Card          Card
	LocationDates map[string][]string
}

type MainData struct {
	Cards        []Card
	CountMembers []int
	Locations    []string
	Members      []string
	CreationDate []string
	FirstAlbum   []string
}

func prepareData(arr []db.Artist) *MainData {
	md := &MainData{
		Cards:        []Card{},
		CountMembers: make([]int, 0),
		Locations:    make([]string, 0),
	}

	temp := make(map[int]int)

	locations := db.DB.GetLocations()
	artist := db.DB.GetArtists()
	for i, artist := range arr {
		num := len(artist.Members)

		temp[num] = 0

		card := Card{
			Id:           artist.Id,
			FirstAlbum:   artist.FirstAlbum,
			Image:        artist.Image,
			GroupName:    artist.Name,
			Location:     locations[i].Location,
			CreationDate: artist.CreationDate,
			Members:      artist.Members,
		}

		md.Cards = append(md.Cards, card)
	}

	for k := range temp {
		md.CountMembers = append(md.CountMembers, k)
	}

	sort.Ints(md.CountMembers)

	uniqueLocations := make(map[string]interface{})
	uniqueMembers := make(map[string]interface{})
	uniqueDate := make(map[string]interface{})
	uniqueFirstAlbum := make(map[string]interface{})

	//----------
	for _, v := range locations {
		for _, v1 := range v.Location {
			uniqueLocations[v1] = nil
		}
	}
	for k := range uniqueLocations {
		md.Locations = append(md.Locations, k)
	}
	//----------
	for _, v := range artist {
		for _, v1 := range v.Members {
			uniqueMembers[v1] = nil
		}
		uniqueDate[strconv.Itoa(v.CreationDate)] = nil

		uniqueDate[v.FirstAlbum] = nil
	}

	for k := range uniqueMembers {
		md.Members = append(md.Members, k)
	}
	for k := range uniqueDate {
		md.CreationDate = append(md.CreationDate, k)
	}
	for k := range uniqueFirstAlbum {
		md.FirstAlbum = append(md.FirstAlbum, k)
	}
	//----------

	sort.Strings(md.Locations)

	return md
}

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Filter struct {
	CreationDateStart int
	CreationDateEnd   int
	AlbumDateStart    int
	AlbumDateEnd      int
	Members           []int
	Location          string
}

func getFilters(r *http.Request) Filter {
	r.ParseForm()

	CdInputStart := r.FormValue("creation-data-from")
	CdInputEnd := r.FormValue("creation-data-to")
	FaInputStart := r.FormValue("first-album-from")
	FaInputEnd := r.FormValue("first-album-to")

	var Members []int

	numMembersValue := r.Form["num-members"]
	for _, numMembers := range numMembersValue {
		numMembersInt, err := strconv.Atoi(numMembers)
		if err == nil {
			Members = append(Members, numMembersInt)
		}
	}

	Location := r.FormValue("location")

	CreationDateStart, _ := strconv.Atoi(CdInputStart)
	CreationDateEnd, _ := strconv.Atoi(CdInputEnd)
	//[len(FaInputStart)-4:]
	AlbumDateStart, _ := strconv.Atoi(FaInputStart)
	AlbumDateEnd, _ := strconv.Atoi(FaInputEnd)

	if CreationDateStart > CreationDateEnd {
		CreationDateStart, CreationDateEnd = CreationDateEnd, CreationDateStart
	}

	if AlbumDateStart > AlbumDateEnd {
		AlbumDateStart, AlbumDateEnd = AlbumDateEnd, AlbumDateStart
	}

	return Filter{
		CreationDateStart,
		CreationDateEnd,
		AlbumDateStart,
		AlbumDateEnd,
		Members,
		Location,
	}
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/html/main-page/index.html", "src/html/main-page/card.html")

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}

	var filteredData []db.Artist

	filter := getFilters(r)

	//fmt.Printf("min: %d, max: %d\n min album: %d, max album: %d\n", startDate, endDate, startAlbumDate, endAlbumDate)

	b := db.DB.GetLocations()
	data := db.DB.GetArtists()

	for i, item := range data {
		a, _ := time.Parse("02-01-2006", item.FirstAlbum)

		year := a.Year()

		if (filter.CreationDateStart != 0 || filter.CreationDateEnd != 0) &&
			(item.CreationDate < filter.CreationDateStart || item.CreationDate > filter.CreationDateEnd) {
			continue
		}

		if (filter.AlbumDateStart != 0 || filter.AlbumDateEnd != 0) &&
			(year < filter.AlbumDateStart || year > filter.AlbumDateEnd) {
			continue
		}

		if len(filter.Members) > 0 && !contains(filter.Members, len(item.Members)) {
			continue
		}

		if filter.Location != "" && !contains(b[i].Location, filter.Location) {
			continue
		}
		filteredData = append(filteredData, item)
	}

	md := prepareData(filteredData)

	err = t.Execute(w, md)

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}
}

func Search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	tmpl, _ := template.ParseFiles("src/html/main-page/index.html", "src/html/main-page/card.html")
	locations := db.DB.GetLocations()
	var filteredData []db.Artist
	search := strings.ToLower(r.FormValue("search-input"))

	for b, artist := range db.DB.GetArtists() {
		if strings.Contains(strings.ToLower(artist.Name), search) {
			filteredData = append(filteredData, artist)
			continue
		}
		if strings.Contains(strings.ToLower(strconv.Itoa(artist.CreationDate)), search) {
			filteredData = append(filteredData, artist)
			continue
		}
		if strings.Contains(strings.ToLower(artist.FirstAlbum), search) {
			filteredData = append(filteredData, artist)
			continue
		}
		for i := range artist.Members {
			if strings.Contains(strings.ToLower(artist.Members[i]), search) {
				filteredData = append(filteredData, artist)
				continue
			}
		}
		for i := range locations[b].Location {
			if strings.Contains(strings.ToLower(locations[b].Location[i]), search) {
				filteredData = append(filteredData, artist)
				continue
			}
		}
	}

	md := prepareData(filteredData)

	tmpl.Execute(w, md)
}

func contains[T string | int](arr []T, compare T) bool {
	for _, v := range arr {
		if v == compare {
			return true
		}
	}

	return false
}
