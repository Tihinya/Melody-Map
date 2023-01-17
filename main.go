package main

import (
	"fmt"
	"groupie-tracker/controllers"
	"groupie-tracker/db"
	"net/http"
)

func main() {
	db := initWithAPIdata()

	http.HandleFunc("/", controllers.MainPage)

	http.ListenAndServe(":8080", nil)

	fmt.Println(db.Dates[:10])
}

// initializes db with API data
func initWithAPIdata() *db.DB {
	var dates db.IndexDates
	var locations db.IndexLocations
	var relations db.IndexRelations
	var artists db.Artists

	db.GetData("https://groupietrackers.herokuapp.com/api/artists", &artists)
	db.GetData("https://groupietrackers.herokuapp.com/api/dates", &dates)
	db.GetData("https://groupietrackers.herokuapp.com/api/locations", &locations)
	db.GetData("https://groupietrackers.herokuapp.com/api/relation", &relations)

	result := db.DB{
		Dates:     dates.Index,
		Locations: locations.Index,
		Relations: relations.Index,
		Artists:   artists,
	}

	return &result
}
