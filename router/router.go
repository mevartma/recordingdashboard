package router

import (
	"RecordingDashboard/db"
	"RecordingDashboard/model"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func NewMux() http.Handler {
	h := http.NewServeMux()
	fs := http.FileServer(http.Dir("templates/"))
	h.Handle("/app/", loggerMid(http.StripPrefix("/app", fs)))
	h.Handle("/api/v1/recordings", loggerMid(http.HandlerFunc(recordingsHandler)))

	return h
}

func recordingsHandler(resp http.ResponseWriter, req *http.Request) {
	var r model.RecordingDetails
	var results []model.RecordingDetails
	var err error
	if req.Method != "GET" {
		err = json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	switch req.Method {
	case "POST":
		err = db.UpdateRecording(r, "add")
	case "GET":
		var idRange model.RecordingSetting
		tmpFrom, err := strconv.Atoi(req.URL.Query().Get("from"))
		tmpTo, err := strconv.Atoi(req.URL.Query().Get("to"))
		idRange.From = int64(tmpFrom)
		idRange.To = int64(tmpTo)

		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		rows, err := db.GetAllRecordingsByRange(idRange.From, idRange.To)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
		}

		for _, rs := range *rows {
			results = append(results, rs)
		}
	default:
		err = errors.New("Method Not Allow")
	}

	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(results)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.Header().Set("Content-type", "application/json")
	resp.Write(js)
	return
}

func loggerMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clIP string
		if r.Header.Get("X-Forwarded-For") == "" {
			clIP = r.RemoteAddr
		} else {
			clIP = r.Header.Get("X-Forwarded-For")
		}
		uAgent := r.Header.Get("User-Agent")
		log.Printf("\"Method\": \"%s\", \"User-Agent\": \"%s\", \"URL\": \"%s\", \"Host\": \"[%s]\", \"Client-IP\": \"%v\"", r.Method, uAgent, r.URL, r.Host, clIP)
		next.ServeHTTP(w, r)
	})
}
