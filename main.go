package main

import (
	"fmt"
	"groupie-tracker/controllers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", controllers.MainPage)
	http.HandleFunc("/full", controllers.FullInfo)
	http.HandleFunc("/search", controllers.Search)
	http.HandleFunc("/filter", controllers.Filter)

	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./src/"))))

	fmt.Println("Your server at: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
