package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/pserson", pserson)
	http.HandleFunc("/form", form)

	http.HandleFunc("/api")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
