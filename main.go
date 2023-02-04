package main

import (
	"groupie-tracker/controllers"
	"groupie-tracker/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router{}

	r.NewRoute("GET", "/", controllers.MainPage)
	r.NewRoute("GET", `/full/(?P<id>\d+)`, controllers.FullInfo)
	r.NewRoute("GET", `/dateslocations/(?P<id>\d+)`, controllers.DatesLocations) // API endpoint for fetching google maps data
	r.NewRoute("GET", `.*`, controllers.NotFound)

	http.HandleFunc("/", r.Serve)

	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./src/"))))

	log.Println("Ctrl + Click on the link: http://localhost:8080")
	log.Println("To stop the server press `Ctrl + C`")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
