package router

import (
	"RecordingDashboard/db"
	"RecordingDashboard/model"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var (
	WorkQueue   = make(chan model.RecordingDetails, 200)
	WorkerQueue chan chan model.RecordingDetails
)

func init() {
	startDispatcher(4)
}

func NewMux() http.Handler {
	h := http.NewServeMux()
	fs := http.FileServer(http.Dir("templates/"))
	h.Handle("/app/", loggerMid(http.StripPrefix("/app", fs)))
	h.Handle("/api/v1/recordings", loggerMid(http.HandlerFunc(recordingsHandler)))

	return h
}

func recordingsHandler(resp http.ResponseWriter, req *http.Request) {
	var results []model.RecordingDetails
	var err error

	switch req.Method {
	case "POST":
		var r model.RecordingDetails
		err = json.NewDecoder(req.Body).Decode(&r)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		//err = db.UpdateRecording(r, "add")
		WorkQueue <- r
	case "GET":
		command := req.URL.Query().Get("command")
		var rows *[]model.RecordingDetails
		if command == "all" {
			rows, err = db.GetAllRecordings()
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else if command == "range" {
			var idRange model.RecordingSetting
			tmpFrom, err := strconv.Atoi(req.URL.Query().Get("from"))
			tmpTo, err := strconv.Atoi(req.URL.Query().Get("to"))
			idRange.From = int64(tmpFrom)
			idRange.To = int64(tmpTo)
			rows, err = db.GetAllRecordingsByRange(idRange.From, idRange.To)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
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

//------------------------------------------Workers------------------------------------------------//

type Worker struct {
	Id          int
	Work        chan model.RecordingDetails
	WorkerQueue chan chan model.RecordingDetails
	QuitChan    chan bool
}

func NewWorker(id int, workerQueue chan chan model.RecordingDetails) Worker {
	worker := Worker{
		Id:          id,
		Work:        make(chan model.RecordingDetails),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
	}
	return worker
}

func (w *Worker) start() {
	go func() {
		for {
			w.WorkerQueue <- w.Work
			select {
			case work := <-w.Work:
				if err := db.UpdateRecording(work, "add"); err != nil {
					log.Println(err)
				}
			case <-w.QuitChan:
				fmt.Printf("Worker id:%d stopping\r\n", w.Id)
				return
			}
		}
	}()
}

func (w *Worker) stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func startDispatcher(nWorkers int) {
	WorkerQueue = make(chan chan model.RecordingDetails, nWorkers)

	for i := 0; i < nWorkers; i++ {
		fmt.Println("Starting Worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				go func() {
					worker := <-WorkerQueue
					worker <- work
				}()
			}
		}
	}()
}

//------------------------------------------Workers------------------------------------------------//
