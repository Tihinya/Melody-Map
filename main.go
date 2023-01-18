package main

import (
	"groupie-tracker/controllers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", controllers.MainPage)
	http.HandleFunc("/full", controllers.FullInfo)

	http.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir("./src/"))))

	log.Println("Ctrl + Click on the link: http://localhost:8080")
	log.Println("To stop the server press `Ctrl + C`")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
