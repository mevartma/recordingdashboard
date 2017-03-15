package utils

import (
	"net/http"
	"RecordingDashboard/model"
	"encoding/json"
	"log"
	"bytes"
	"encoding/hex"
	"time"
	"fmt"
	"crypto/sha1"
)

var httpClient = http.Client{}


func ValidateUserName(u model.UserNameAndPassword, g string) (bool, error) {
	userURL := "http://192.168.1.66/api/phonebook/it/validate"
	var result model.UserResult

	js, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
		return false,err
	}

	req, err := http.NewRequest("POST", userURL, bytes.NewReader(js))
	if err != nil {
		log.Fatal(err)
		return false,err
	}
	req.Header.Set("Content-Type", "application/json")
	respUserURL, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return false,err
	}

	err = json.NewDecoder(respUserURL.Body).Decode(&result)
	if err != nil {
		log.Fatal(err)
		return false,err
	}

	if result.Username == u.Username && result.Result == false {
		return false, nil
	} else if result.Username == u.Username && result.Result == true {
		groupURL := "http://192.168.1.66/api/phonebook/it/usergroups"
		var username model.UserName
		var usergroups []model.UserGroups
		username.Username = u.Username

		js, err = json.Marshal(username)
		if err != nil {
			log.Fatal(err)
			return false,err
		}

		req, err := http.NewRequest("POST", groupURL, bytes.NewReader(js))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		respGroupURL, err := httpClient.Do(req)
		if err != nil {
			log.Fatal(err)
			return false,err
		}

		err = json.NewDecoder(respGroupURL.Body).Decode(&usergroups)
		if err != nil {
			log.Fatal(err)
			return false,err
		}
		t := false
		for _, group := range usergroups {
			if group.GroupName == g  {
				t = true
			}
		}
		return t, nil
	}

	return false,err
}

func CreateSessionCoockie(toHash string, t time.Time) string {
	h := sha1.New()
	stringToHash := fmt.Sprintf("%s%v",toHash,t)
	h.Write([]byte(stringToHash))

	sha512String := hex.EncodeToString(h.Sum(nil))
	return sha512String
}