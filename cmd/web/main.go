package main

import (
	"log"
	"net/http"
)

func main() {
	router := GetRouter()
	log.Println("Listenting on 127.0.0.1:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
