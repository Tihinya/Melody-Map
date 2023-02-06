package controllers

import (
	"groupie-tracker/db"
	"net/http"
	"strconv"
	"strings"
)

type Filter struct {
	CreationDateStart int
	CreationDateEnd   int
	AlbumDateStart    int
	AlbumDateEnd      int
	Members           []int
	Location          string
}

func getFilters(r *http.Request) Filter {
	r.ParseForm()

	CdInputStart := r.FormValue("creation-data-from")
	CdInputEnd := r.FormValue("creation-data-to")
	FaInputStart := r.FormValue("first-album-from")
	FaInputEnd := r.FormValue("first-album-to")

	var Members []int

	numMembersValue := r.Form["num-members"]
	for _, numMembers := range numMembersValue {
		numMembersInt, err := strconv.Atoi(numMembers)
		if err == nil {
			Members = append(Members, numMembersInt)
		}
	}

	Location := r.FormValue("location")

	CreationDateStart, _ := strconv.Atoi(CdInputStart)
	CreationDateEnd, _ := strconv.Atoi(CdInputEnd)
	AlbumDateStart, _ := strconv.Atoi(FaInputStart)
	AlbumDateEnd, _ := strconv.Atoi(FaInputEnd)

	// errors escaped because later on
	if CreationDateStart > CreationDateEnd {
		CreationDateStart, CreationDateEnd = CreationDateEnd, CreationDateStart
	}

	if AlbumDateStart > AlbumDateEnd {
		AlbumDateStart, AlbumDateEnd = AlbumDateEnd, AlbumDateStart
	}

	return Filter{
		CreationDateStart,
		CreationDateEnd,
		AlbumDateStart,
		AlbumDateEnd,
		Members,
		Location,
	}
}

func contains[T string | int](arr []T, compare T) bool {
	for _, v := range arr {
		if v == compare {
			return true
		}
	}

	return false
}

func partialContains(arr []string, s string) bool {
	for _, v := range arr {
		if strings.Contains(strings.ToLower(v), s) {
			return true
		}
	}

	return false
}

func isArtist(artist db.Artist, search string, locations []string) bool {
	return strings.Contains(strings.ToLower(artist.Name), search) ||
		strings.Contains(strings.ToLower(strconv.Itoa(artist.CreationDate)), search) ||
		strings.Contains(strings.ToLower(artist.FirstAlbum), search) ||
		partialContains(artist.Members, search) ||
		partialContains(locations, search)
}

func mapToSlice[T int | string](m map[T]interface{}) []T {
	arr := []T{}

	for k := range m {
		arr = append(arr, k)
	}

	return arr
}
