package main

import (
	"RecordingDashboard/router"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", router.NewMux()))
}