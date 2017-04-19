package main

import (
	"RecordingDashboard/router"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":80", router.NewMux()))
}
