package router

import (
	"RecordingDashboard/db"
	"RecordingDashboard/model"
	"RecordingDashboard/utils"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	h.Handle("/app", loggerMid(authMid(http.HandlerFunc(appPage))))
	h.Handle("/api/v1/recordings", loggerMid(authMid(http.HandlerFunc(recordingsHandler))))
	h.Handle("/login", loggerMid(http.HandlerFunc(loginPage)))
	h.Handle("/api/v1/users/loginuser", loggerMid(http.HandlerFunc(usersLoginHandler)))
	h.Handle("/api/v1/users/logoutuser", loggerMid(http.HandlerFunc(usersLogoutHandler)))
	h.Handle("/favicon.ico", loggerMid(http.HandlerFunc(fav)))
	h.Handle("/", loggerMid(http.HandlerFunc(home)))

	return h
}

func fav(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
}

func loginPage(resp http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/login.html")
	t.Execute(resp, nil)
}

func appPage(resp http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(resp,nil)
}

func home(resp http.ResponseWriter, req *http.Request) {
	http.Redirect(resp, req, "/login", http.StatusFound)
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

func usersLoginHandler(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	var user model.UserNameAndPassword
	var realUser model.UserDetails

	user.Username = strings.Replace(fmt.Sprintf("%s", req.Form["username"]), "[", "", -1)
	user.Username = strings.Replace(user.Username, "]", "", -1)
	user.Password = strings.Replace(fmt.Sprintf("%s", req.Form["password"]), "[", "", -1)
	user.Password = strings.Replace(user.Password, "]", "", -1)

	result, err := utils.ValidateUserName(user, "ITGroup")
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	if result == false {
		http.Redirect(resp, req, "/login", http.StatusFound)
		return
	}

	exprDate := time.Now().AddDate(0, 0, 1)

	sessionID := utils.CreateSessionCoockie(user.Username, exprDate)
	var clIP string
	if req.Header.Get("X-Forwarded-For") == "" {
		clIP = req.RemoteAddr
	} else {
		clIP = req.Header.Get("X-Forwarded-For")
	}

	realUser.UserName = user.Username
	realUser.IpAddress = clIP
	realUser.UserAgent = req.Header.Get("User-Agent")
	realUser.Cookie = sessionID
	realUser.ExpireTime = exprDate.String()
	err = db.UpdateUser(realUser, "add")
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookieMonster := &http.Cookie{
		Name:     "SessionID",
		Expires:  exprDate,
		Value:    sessionID,
		HttpOnly: true,
		MaxAge:   50000,
		Path:     "/",
	}

	http.SetCookie(resp, cookieMonster)
	http.Redirect(resp, req, "/app", http.StatusFound)
	return
}

func usersLogoutHandler(resp http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("SessionID")
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = db.DeleteSessionId(cookie.Value)

	exprDate := time.Now()
	cookieMonster := &http.Cookie{
		Name:    "SessionID",
		Expires: exprDate,
		Value:   "",
	}

	http.SetCookie(resp, cookieMonster)
	http.Redirect(resp, req, "/login", http.StatusFound)
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

func authMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clIP string
		if r.Header.Get("X-Forwarded-For") == "" {
			clIP = r.RemoteAddr
		} else {
			clIP = r.Header.Get("X-Forwarded-For")
		}
		uAgent := r.Header.Get("User-Agent")

		cookie, err := r.Cookie("SessionID")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if cookie == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		status, user, err := db.GetSessionId(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if status == true && user.UserAgent == uAgent && user.IpAddress == clIP {
			next.ServeHTTP(w, r)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
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
