package controllers

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/db"
	"groupie-tracker/router"
	"html/template"
	"log"
	"net/http"
	"os"
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
}

func prepareData(arr []db.Artist) *MainData {
	md := &MainData{
		Cards:        []Card{},
		CountMembers: make([]int, 0),
		Locations:    make([]string, 0),
	}

	temp := make(map[int]int)

	for _, artist := range arr {
		num := len(artist.Members)

		temp[num] = 0

		card := Card{
			Id:           artist.Id,
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

	locat := make(map[string]interface{})

	locations := db.DB.GetLocations()

	for _, v := range locations {
		for _, v1 := range v.Location {
			locat[v1] = nil
		}

	}

	for k := range locat {
		md.Locations = append(md.Locations, k)
	}

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

		if item.CreationDate < filter.CreationDateStart || item.CreationDate > filter.CreationDateEnd {
			continue
		}

		if year < filter.AlbumDateStart || year > filter.AlbumDateEnd {
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

// func Search(w http.ResponseWriter, r *http.Request) {
// 	groupName := r.FormValue("search-input")

// 	if r.Method != "GET" {
// 		fmt.Println(http.StatusBadRequest)
// 		return
// 	}

// 	http.RedirectHandler("http://localhost:8080/full/2", http.StatusMovedPermanently)

// 	fmt.Println(groupName)
// }

func contains[T string | int](arr []T, compare T) bool {
	for _, v := range arr {
		if v == compare {
			return true
		}
	}

	return false
}

// func Filter(w http.ResponseWriter, r *http.Request) {
// 	var filteredData []db.Artist

// 	r.ParseForm()
// 	if r.Method != "GET" {
// 		w.WriteHeader(http.StatusBadRequest)
// 	}
// 	CdInputStart := r.FormValue("creation-data-from")
// 	CdInputEnd := r.FormValue("creation-data-to")
// 	FaInputStart := r.FormValue("first-album-from")
// 	FaInputEnd := r.FormValue("first-album-to")

// 	var memberSlice []int

// 	numMembersValue := r.Form["num-members"]
// 	for _, numMembers := range numMembersValue {
// 		numMembersInt, err := strconv.Atoi(numMembers)
// 		if err == nil {
// 			memberSlice = append(memberSlice, numMembersInt)
// 		}
// 	}

// 	locationInput := r.FormValue("location")

// 	data := db.DB.GetArtists()

// 	startDate, _ := strconv.Atoi(CdInputStart)
// 	endDate, _ := strconv.Atoi(CdInputEnd)
// 	//[len(FaInputStart)-4:]
// 	startAlbumDate, _ := strconv.Atoi(FaInputStart)
// 	endAlbumDate, _ := strconv.Atoi(FaInputEnd)

// 	if startDate > endDate {
// 		startDate, endDate = endDate, startDate
// 	}

// 	if startAlbumDate > endAlbumDate {
// 		startAlbumDate, endAlbumDate = endAlbumDate, startAlbumDate
// 	}

// 	//fmt.Printf("min: %d, max: %d\n min album: %d, max album: %d\n", startDate, endDate, startAlbumDate, endAlbumDate)

// 	b := db.DB.GetLocations()

// 	for i, item := range data {
// 		a, _ := time.Parse("02-01-2006", item.FirstAlbum)

// 		year := a.Year()

// 		if item.CreationDate < startDate || item.CreationDate > endDate {
// 			continue
// 		}

// 		if year < startAlbumDate || year > endAlbumDate {
// 			continue
// 		}

// 		if len(memberSlice) > 0 && !contains(memberSlice, len(item.Members)) {
// 			continue
// 		}

// 		if locationInput != "" && !contains(b[i].Location, locationInput) {
// 			continue
// 		}
// 		filteredData = append(filteredData, item)
// 	}
// 	md := prepareData(filteredData)

// 	t, err := template.ParseFiles("src/html/main-page/index.html", "src/html/main-page/card.html")
// 	if err != nil {
// 		fmt.Println(http.StatusInternalServerError, err)
// 	}

// 	err = t.Execute(w, md)
// 	if err != nil {
// 		fmt.Println(http.StatusInternalServerError, err)
// 	}
// }

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
			res, err := http.Get(req)
			if err != nil {
				log.Println(err)
			}

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
