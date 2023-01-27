package db

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type googleMap struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func GetGoogleMap(r *http.Response) (lat float64, lng float64) {
	body, err := io.ReadAll(r.Body)

	var gm googleMap

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &gm)

	if err != nil {
		log.Fatal(err)
	}

	loc := gm.Results[0].Geometry.Location

	return loc.Lat, loc.Lng
}
