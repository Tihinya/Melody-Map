package main

import (
	"fmt"
	"groupie-tracker/controllers"
	"net/http"
)

func main() {
	http.HandleFunc("/", controllers.MainPage)
	http.HandleFunc("/full", controllers.FullInfo)
	http.HandleFunc("/search", controllers.Search)

	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./src/"))))

	fmt.Println("Your server at: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
