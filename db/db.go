package db

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// API schemas
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

type IndexLocations struct {
	Index Locations `json:"index"`
}

type Date struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Dates []Date

type IndexDates struct {
	Index Dates `json:"index"`
}

type Relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
type Relations []Relation

type IndexRelations struct {
	Index Relations `json:"index"`
}

type DB struct {
	Dates     Dates
	Locations Locations
	Relations Relations
	Artists   Artists
}

type apiData interface {
	*IndexDates | *IndexLocations | *IndexRelations | *Artists
}

func GetData[T apiData](url string, schema T) {
	r, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &schema)

	if err != nil {
		log.Fatal(err)
	}
}
