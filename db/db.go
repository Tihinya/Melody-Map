package db

import (
	"encoding/json"
	"io"
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

type database struct {
	Dates     Dates
	Locations Locations
	Relations Relations
	Artists   Artists
}

type apiData interface {
	*IndexDates | *IndexLocations | *IndexRelations | *Artists
}

var DB *database

func GetData[T apiData](url string, schema T) {
	r, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &schema)

	if err != nil {
		log.Fatal(err)
	}
}

// TODO takes a custom filter
func (db *database) GetArtists() Artists {
	if db.Artists != nil {
		return db.Artists
	}

	return nil
}

func (db *database) GetLocations() Locations {
	if db.Locations != nil {
		return db.Locations
	}

	return nil
}

func (db *database) GetRelations() Relations {
	if db.Relations != nil {
		return db.Relations
	}

	return nil
}

func (db *database) GetDates() Dates {
	if db.Dates != nil {
		return db.Dates
	}

	return nil
}

func (db *database) GetAllRecords() {}

func (db *database) GetArtistById() {}

// initializes db with API data
func initWithAPIdata() *database {
	var dates IndexDates
	var locations IndexLocations
	var relations IndexRelations
	var artists Artists

	GetData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	GetData("https://groupietrackers.herokuapp.com/api/dates", &dates)
	GetData("https://groupietrackers.herokuapp.com/api/locations", &locations)
	GetData("https://groupietrackers.herokuapp.com/api/relation", &relations)

	result := database{
		Dates:     dates.Index,
		Locations: locations.Index,
		Relations: relations.Index,
		Artists:   artists,
	}

	return &result
}

func init() {
	DB = initWithAPIdata()
}
