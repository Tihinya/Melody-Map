package db

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// API schemas
type artist struct {
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

type artists []artist

type location struct {
	Id       int      `json:"id"`
	Location []string `json:"locations"`
	Dates    string   `json:"dates"`
}

type locations []location

type indexLocations struct {
	Index locations `json:"index"`
}

type date struct {
	Id    int      `json:"id"`
	Dates []string `json:"date"`
}

type dates []date

type indexDates struct {
	Index dates `json:"index"`
}

type relation struct {
	Id             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}
type relations []relation

type indexRelations struct {
	Index relations `json:"index"`
}

type database struct {
	Dates     dates
	Locations locations
	Relations relations
	Artists   artists
}

type apiData interface {
	*indexDates | *indexLocations | *indexRelations | *artists
}

var DB database

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
func (db *database) GetArtists() artists {
	if db.Artists != nil {
		return db.Artists
	}

	return nil
}

func (db *database) GetLocations() locations {
	if db.Locations != nil {
		return db.Locations
	}

	return nil
}
func (db *database) GetRelations() relations {
	if db.Relations != nil {
		return db.Relations
	}
	return nil
}
func (db *database) GetDates() dates {
	if db.Dates != nil {
		return db.Dates
	}

	return nil
}

func (db *database) GetAllRecords() {}

func (db *database) GetArtistById() {}

// initializes db with API data
func initWithAPIdata() *database {
	var dates indexDates
	var locations indexLocations
	var relations indexRelations
	var artists artists

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
	DB = *initWithAPIdata()
}
