package controllers

import (
	"fmt"
	"groupie-tracker/db"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"
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

	locat := make(map[string]int)

	locations := db.DB.GetLocations()

	for i := range md.Cards {
		md.Cards[i].Location = locations[i].Location
		var loc string
		for _, a := range locations[i].Location {
			loc = a
		}
		locat[loc] = i
	}

	for k := range locat {
		md.Locations = append(md.Locations, k)
	}

	return md
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/html/main-page/index.html", "src/html/main-page/card.html")

	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}

	md := prepareData(db.DB.GetArtists())

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
	groupName := r.FormValue("search-input")

	if r.Method != "GET" {
		fmt.Println(http.StatusBadRequest)
		return
	}

	fmt.Println(groupName)

}

func Filter(w http.ResponseWriter, r *http.Request) {
	var filteredData []db.Artist

	r.ParseForm()
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	}
	CdInputStart := r.FormValue("creation-data-from")
	CdInputEnd := r.FormValue("creation-data-to")
	FaInputStart := r.FormValue("first-album-from")
	FaInputEnd := r.FormValue("first-album-to")

	var memberSlice []int

	numMembersValue := r.Form["num-members"]
	for _, numMembers := range numMembersValue {
		numMembersInt, err := strconv.Atoi(numMembers)
		if err == nil {
			memberSlice = append(memberSlice, numMembersInt)
		}
	}

	Location := r.FormValue("location")

	data := db.DB.GetArtists()

	startDate, _ := strconv.Atoi(CdInputStart)
	endDate, _ := strconv.Atoi(CdInputEnd)
	//[len(FaInputStart)-4:]
	startAlbumDate, _ := strconv.Atoi(FaInputStart)
	endAlbumDate, _ := strconv.Atoi(FaInputEnd)

	if startDate > endDate {
		startDate, endDate = endDate, startDate
	}

	if startAlbumDate > endAlbumDate {
		startAlbumDate, endAlbumDate = endAlbumDate, startAlbumDate
	}

	//fmt.Printf("min: %d, max: %d\n min album: %d, max album: %d\n", startDate, endDate, startAlbumDate, endAlbumDate)

	for _, item := range data {
		a, _ := time.Parse("02-01-2006", item.FirstAlbum)

		year := a.Year()

		if item.CreationDate < startDate || item.CreationDate > endDate {
			continue
		}
		if year < startAlbumDate || year > endAlbumDate {
			continue
		}
		if len(memberSlice) == 0 && Location == "" {
			continue
		} else if len(memberSlice) != 0 && Location == "" {
			for _, v := range memberSlice {
				if len(item.Members) == v {
					filteredData = append(filteredData, item)
				}
			}
		} else if len(memberSlice) == 0 && Location != "" {
			for _, location := range db.DB.Locations {
				if item.Id == location.Id {
					for _, place := range location.Location {
						if Location == place {
							filteredData = append(filteredData, item)
						}
					}
				}
			}
		} else {
			for _, v := range memberSlice {
				if len(item.Members) == v {
					for _, location := range db.DB.Locations {
						if item.Id == location.Id {
							for _, place := range location.Location {
								if Location == place {
									filteredData = append(filteredData, item)
								}
							}

						}
					}
				}
			}
		}
	}
	fmt.Println(filteredData)
	md := prepareData(filteredData)

	t, err := template.ParseFiles("src/html/main-page/index.html", "src/html/main-page/card.html")
	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}

	err = t.Execute(w, md)
	if err != nil {
		fmt.Println(http.StatusInternalServerError, err)
	}

}
