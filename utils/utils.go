package utils

import (
	"net/http"
	"RecordingDashboard/model"
	"encoding/json"
	"log"
	"bytes"
	"crypto/sha512"
	"encoding/hex"
)

var client = &http.Client{}


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
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
		return false,err
	}

	err = json.NewDecoder(req.Body).Decode(&result)
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
		_, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
			return false,err
		}

		err = json.NewDecoder(req.Body).Decode(&usergroups)
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

func CreateSessionCoockie(toHash string) string {
	h := sha512.New()
	h.Write([]byte(toHash))

	sha512String := hex.EncodeToString(h.Sum(nil))
	return sha512String
}