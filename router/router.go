package router

import (
	"RecordingDashboard/db"
	"RecordingDashboard/model"
	"RecordingDashboard/utils"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func NewMux() http.Handler {
	h := http.NewServeMux()
	h.Handle("/app", loggerMid(authMid(http.HandlerFunc(appPage))))
	h.Handle("/api/v1/recordings", loggerMid(authMid(http.HandlerFunc(recordingsHandler))))
	h.Handle("/login", loggerMid(http.HandlerFunc(loginPage)))
	h.Handle("/api/v1/users/loginuser", loggerMid(http.HandlerFunc(usersLoginHandler)))
	h.Handle("/api/v1/users/logoutuser", loggerMid(http.HandlerFunc(usersLogoutHandler)))
	h.Handle("/favicon.ico", loggerMid(http.HandlerFunc(fav)))
	h.Handle("/", loggerMid(http.HandlerFunc(files)))

	return h
}

func fav(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
}

func loginPage(resp http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(resp, "login.html", nil)
}

func appPage(resp http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(resp, "index.html", nil)
}

func files(resp http.ResponseWriter, req *http.Request) {
	var data []byte
	var err error
	var fileType string
	if req.URL.String() == "/" {
		http.Redirect(resp, req, "/login", http.StatusFound)
	} else if strings.Contains(req.URL.String(), "betamediarecording") {
		tempURL := fmt.Sprintf("https://s3.eu-central-1.amazonaws.com%s", req.URL.String())
		amazonURL := strings.Replace(tempURL, "mp3", "gsm", -1)

		splitFolder := strings.Split(req.URL.String(), "/")
		folderPath := fmt.Sprintf("temp%s/%s", splitFolder[0], splitFolder[1])

		if _, err = os.Stat(folderPath); os.IsNotExist(err) {
			os.MkdirAll(folderPath, os.ModePerm)
		}

		check := http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				req.URL.Opaque = req.URL.Path
				return nil
			},
		}

		response, err := check.Get(amazonURL)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()
		fileName := strings.Replace(splitFolder[len(splitFolder)-1], "mp3", "gsm", -1)
		newFileName := fmt.Sprintf("%s/%s", folderPath, splitFolder[len(splitFolder)-1])

		filePath := fmt.Sprintf("%s/%s", folderPath, fileName)
		out, err := os.Create(filePath)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(out, response.Body)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		out.Close()

		cmdName := "sox"
		ConvertCommand := exec.Command(cmdName, filePath, "-S", newFileName)
		_, err = ConvertCommand.CombinedOutput()
		if err != nil {
			fmt.Println(err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		data, err = ioutil.ReadFile(newFileName)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		fileType = "mp3"

		/*os.Remove(filePath)
		os.Remove(newFileName)*/
	} else {
		filePath := fmt.Sprintf("templates%s", req.URL.String())
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}
		fileExt := strings.Split(filePath, ".")
		fileType = fileExt[len(fileExt)-1]
	}

	switch fileType {
	case "js":
		resp.Header().Set("Content-Type", "application/javascript")
	case "css":
		resp.Header().Set("Content-Type", "text/css")
	case "mp3":
		resp.Header().Set("Content-Type", "audio/mpeg;audio/mpeg3;audio/x-mpeg-3;video/mpeg;video/x-mpeg;text/xml")
	}

	resp.Write(data)
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
	case "GET":
		command := req.URL.Query().Get("command")
		var rows *[]model.RecordingDetails
		switch command {
		case "number":
			number := req.URL.Query().Get("number")
			date1 := req.URL.Query().Get("date1")
			date2 := req.URL.Query().Get("date1")
			office := req.URL.Query().Get("office")
			if number == "" || date1 == "" || date2 == "" || office == "" {
				http.Error(resp, "missing details", http.StatusBadRequest)
			}
			rows, err = db.GetRecording(number, date1, date2, office)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			http.Error(resp, "Command not found", http.StatusUnauthorized)
			return
		}

		for _, rs := range *rows {
			results = append(results, rs)
		}
	default:
		http.Error(resp, "Command not found", http.StatusUnauthorized)
		return
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

	result, err := utils.ValidateUserName(user, "RecrodingSystem")
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

	realIP := strings.Split(clIP, ":")
	clIP = realIP[0]

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

	cookieMonster := &http.Cookie{
		Name:   "SessionID",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(resp, cookieMonster)
	http.Redirect(resp, req, "/login", http.StatusFound)
}

func loggerMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		req.Close = true
		var clIP string
		if req.Header.Get("X-Forwarded-For") == "" {
			clIP = req.RemoteAddr
		} else {
			clIP = req.Header.Get("X-Forwarded-For")
		}
		realip := strings.Split(clIP, ":")
		clIP = realip[0]
		uAgent := req.Header.Get("User-Agent")
		log.Printf("\"Method\": \"%s\", \"User-Agent\": \"%s\", \"URL\": \"%s\", \"Host\": \"[%s]\", \"Client-IP\": \"%v\"", req.Method, uAgent, req.URL, req.Host, clIP)
		next.ServeHTTP(resp, req)
	})
}

func authMid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		var clIP string
		if req.Header.Get("X-Forwarded-For") == "" {
			clIP = req.RemoteAddr
		} else {
			clIP = req.Header.Get("X-Forwarded-For")
		}

		realip := strings.Split(clIP, ":")
		clIP = realip[0]

		if clIP == "192.168.50.14" || clIP == "192.168.150.113" {
			next.ServeHTTP(resp, req)
		} else {
			uAgent := req.Header.Get("User-Agent")

			cookie, err := req.Cookie("SessionID")
			if err != nil {
				http.Redirect(resp, req, "/login", http.StatusFound)
				return
			}

			if cookie == nil {
				http.Redirect(resp, req, "/login", http.StatusFound)
				return
			}

			status, user, err := db.GetSessionId(cookie.Value)
			if err != nil {
				http.Redirect(resp, req, "/login", http.StatusFound)
				return
			}

			if status == true && user.UserAgent == uAgent && user.IpAddress == clIP {
				next.ServeHTTP(resp, req)
			} else {
				http.Redirect(resp, req, "/login", http.StatusFound)
				return
			}
		}
	})
}
