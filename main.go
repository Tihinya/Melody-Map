package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Artist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CrationDate  int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Artists []Artist

type Location struct {
	Id       int      `json:"id"`
	Location []string `json:"locations"`
	Dates    string   `json:"dates"`
}

type Locations []Location

type IndexOfLocations struct {
	Index Locations `json:"index"`
}

func main() {
	var arts Artists

	r, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")

	if err != nil {
		return
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &arts)

	if err != nil {
		fmt.Print("a")
	}

	fmt.Println(arts[0])
	loct()
}

func loct() {
	var locations IndexOfLocations

	r, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")

	if err != nil {
		return
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Print(err)
	}

	err = json.Unmarshal(body, &locations)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(locations.Index[0].Id)
}
