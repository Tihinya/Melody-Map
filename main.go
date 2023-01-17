package main

import (
	"groupie-tracker/controllers"
	"net/http"
)

func main() {
	http.HandleFunc("/", controllers.MainPage)
	http.HandleFunc("/full", controllers.FullInfo)

	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./src/"))))

	http.ListenAndServe(":8080", nil)
}
