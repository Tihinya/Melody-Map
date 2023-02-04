package controllers

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/db"
	"groupie-tracker/errorsSafe"
	"groupie-tracker/router"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func DatesLocations(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	k := os.Getenv("GOOGLE_API_KEY")

	sid := router.GetField(r, "id")

	id, err := strconv.Atoi(sid)

	if err != nil {
		errorsSafe.WrapError(err, errorsSafe.ErrServer)
		return
	}

	var coordinates []Coordinate

	if k == "" {
		log.Println("Need to enter google map API key to GOOGLE_API_KEY env variable")
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

			lat, lng, err := db.GetGoogleMap(res)
			if err != nil {
				log.Println(err)
				continue
			}

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
