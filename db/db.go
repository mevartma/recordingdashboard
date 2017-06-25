package db

import (
	"RecordingDashboard/model"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const (
	db_local_URL   string = "root:@tcp(localhost:3306)/recordingdb"
	db_germany_URL string = "recording:recording@tcp(192.168.50.14:3306)/asteriskcdrdb"
	db_kiev_URL    string = "recording:recording@tcp(localhost:3306)/asteriskcdrdb"
	server         string = "mysql"
)

var sqlDatabase *sql.DB
var err error

func GetSessionId(s string) (bool, *model.UserDetails, error) {
	err = nil
	var results model.UserDetails
	sqlDatabase, err = sql.Open(server, db_local_URL)
	defer sqlDatabase.Close()
	query := "SELECT * FROM userssessions WHERE cookie = ? LIMIT 1"
	rows, err := sqlDatabase.Query(query, s)
	if err != nil {
		return false, nil, err
	}
	for rows.Next() {
		err = rows.Scan(&results.Id, &results.UserName, &results.IpAddress, &results.UserAgent, &results.Cookie, &results.ExpireTime)
		if err != nil {
			return false, nil, err
		}
	}

	if results.UserName != "" {
		return true, &results, err
	}

	return false, &results, err
}

func GetAllSessions() (*[]model.UserDetails, error) {
	err = nil
	var result []model.UserDetails
	sqlDatabase, err = sql.Open(server, db_local_URL)
	defer sqlDatabase.Close()
	query := "SELECT * FROM userssessions"
	rows, err := sqlDatabase.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.UserDetails
		err = rows.Scan(&r.Id, &r.UserName, &r.IpAddress, &r.UserAgent, &r.Cookie, &r.ExpireTime)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
	}

	return &result, err
}

func DeleteSessionId(s string) error {
	err = nil
	sqlDatabase, err = sql.Open(server, db_local_URL)
	if err != nil {
		return err
	}
	defer sqlDatabase.Close()

	stmt, err := sqlDatabase.Prepare("DELETE FROM userssessions WHERE cookie = ?")
	_, err = stmt.Exec(s)

	return err
}

func UpdateUser(r model.UserDetails, action string) error {
	err = nil
	sqlDatabase, err = sql.Open(server, db_local_URL)
	defer sqlDatabase.Close()

	switch action {
	case "add":
		stmt, err := sqlDatabase.Prepare("INSERT INTO userssessions(username,ipaddress,useragent,cookie,expiretime) VALUES (?,?,?,?,?)")
		_, err = stmt.Exec(r.UserName, r.IpAddress, r.UserAgent, r.Cookie, r.ExpireTime)
		return err
	default:
		err = errors.New("Command Not Found")
	}

	return err
}

func GetRecording(num, date1, date2, office string) (*[]model.RecordingDetails, error) {
	var results []model.RecordingDetails
	err = nil

	switch office {
	case "germany":
		sqlDatabase, err = sql.Open(server, db_germany_URL)
	case "'kiev":
		sqlDatabase, err = sql.Open(server, db_kiev_URL)
	default:
		err = errors.New("Choose Right Office")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer sqlDatabase.Close()

	query := "SELECT calldate,clid,src,dst,duration,billsec,disposition,cnam,recordingfile,uniqueid FROM cdr WHERE calldate BETWEEN ? AND ? AND disposition like 'ANSWERED' AND  ( src LIKE ? OR dst LIKE ? )"
	rows, err := sqlDatabase.Query(query, date1, date2, num, num)
	c, _ := rows.Columns()
	fmt.Println(c)
	for rows.Next() {
		var r model.RecordingDetails
		err = rows.Scan(&r.CallDate, &r.ClId, &r.SRC, &r.DST, &r.Duration, &r.BillSec, &r.Disposition, &r.Cnam, &r.RecordingFile, &r.UniqueId)
		fmt.Println(r)
		r.Office = office
		switch office {
		case "germany":
			r.FileURL = fmt.Sprintf("https://s3.eu-central-1.amazonaws.com/betamediarecording/Germany/%s", r.RecordingFile)
			r.ServerURL = fmt.Sprintf("http://192.168.1.7/betamediarecording/Germany/%s", r.RecordingFile)
		case "kiev":
			r.FileURL = fmt.Sprintf("https://s3.eu-central-1.amazonaws.com/betamediarecording/Kiev/%s", r.RecordingFile)
			r.ServerURL = fmt.Sprintf("http://192.168.1.7/betamediarecording/Kiev/%s", r.RecordingFile)
		}

		results = append(results, r)
	}

	return &results, err
}
