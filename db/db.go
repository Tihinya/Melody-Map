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
	CreationDate int      `json:"creationDate"`
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

// type IDB interface {
// 	*Dates | *Locations | *Relations | *Artists
// }

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

// TODO takes a custom filter
func (db *DB) GetArtists() Artists {
	if db.Artists != nil {
		return db.Artists
	}

	return nil
}

func (db *DB) GetLocations() Locations {
	if db.Locations != nil {
		return db.Locations
	}

	return nil
}
func (db *DB) GetRelations() Relations {
	if db.Relations != nil {
		return db.Relations
	}

	return nil
}
func (db *DB) GetDates() Dates {
	if db.Dates != nil {
		return db.Dates
	}

	return nil
}

func (db *DB) GetAllRecords() {}

func (db *DB) GetArtistById() {}

// return every match on key
// func (db *DB) Search(v string) []string {
// 	var result []string

// 	for {
// 		// exit case on mathc
// 	}
// }
