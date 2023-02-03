package db

import (
	"encoding/json"
	"fmt"
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

func GetGoogleMap(r *http.Response) (lat float64, lng float64, err error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var gm googleMap
	err = json.Unmarshal(body, &gm)

	if err != nil {
		log.Println(err)
		return
	}

	if len(gm.Results) < 1 {
		return 0, 0, fmt.Errorf("cannot get latitude, longtitude")
	}
	loc := gm.Results[0].Geometry.Location

	return loc.Lat, loc.Lng, err
}
