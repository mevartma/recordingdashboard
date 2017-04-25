package db

import (
	"RecordingDashboard/model"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

const (
	dbURL  string = "root:@tcp(localhost:3306)/recordingdb"
	server string = "mysql"
)

func UpdateRecording(r model.RecordingDetails, action string) error {
	var err error
	db, err := sql.Open(server, dbURL)
	defer db.Close()

	switch action {
	case "add":
		stmt, err := db.Prepare("INSERT INTO recordings(calldate,clid,src,dst,duration,billsec,disposition,accountcode,uniqueid,did,recordingfile,diskfilepath,s3fileurl,office) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
		_, err = stmt.Exec(r.CallDate, r.ClId, r.SRC, r.DST, r.Duration, r.BillSec, r.Disposition, r.AccountCode, r.UniqueId, r.DID, r.RecordingFile, r.DiskFilePath, r.S3FileURL, r.Office)
		return err
	default:
		err = errors.New("Command Not Found")
	}

	return err
}

func GetSessionId(s string) (bool, *model.UserDetails, error) {
	var err error
	var results model.UserDetails
	db, err := sql.Open(server, dbURL)
	defer db.Close()
	query := "SELECT * FROM userssessions WHERE cookie = ? LIMIT 1"
	rows, err := db.Query(query, s)
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

func DeleteSessionId(s string) error {
	db, err := sql.Open(server, dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM userssessions WHERE cookie = ?")
	_, err = stmt.Exec(s)

	return err
}

func UpdateUser(r model.UserDetails, action string) error {

	var err error
	db, err := sql.Open(server, dbURL)
	defer db.Close()

	switch action {
	case "add":
		stmt, err := db.Prepare("INSERT INTO userssessions(username,ipaddress,useragent,cookie,expiretime) VALUES (?,?,?,?,?)")
		_, err = stmt.Exec(r.UserName, r.IpAddress, r.UserAgent, r.Cookie, r.ExpireTime)
		return err
	default:
		err = errors.New("Command Not Found")
	}

	return err
}

func GetAllRecordings() (*[]model.RecordingDetails, error) {
	var results []model.RecordingDetails
	db, err := sql.Open(server, dbURL)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id,calldate,src,dst,duration,billsec,disposition,s3fileurl,office FROM recordings"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.RecordingDetails
		err = rows.Scan(&r.Id, &r.CallDate, &r.SRC, &r.DST, &r.Duration, &r.BillSec, &r.Disposition, &r.S3FileURL, &r.Office)
		if err != nil {
			return nil, err
		}
		tempURL := strings.Replace(r.S3FileURL,"https://s3.eu-central-1.amazonaws.com/","http://192.168.1.7/",-1)
		tempURL2 := strings.Replace(tempURL,"gsm","wav",-1)
		r.S3FileURL = tempURL2
		results = append(results, r)
	}

	return &results, err
}

func GetRecordingsByRange(from, to int64) (*[]model.RecordingDetails, error) {
	var results []model.RecordingDetails
	db, err := sql.Open(server, dbURL)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id,calldate,src,dst,duration,billsec,disposition,s3fileurl,office FROM recordings WHERE id BETWEEN ? AND ?"
	rows, err := db.Query(query, from, to)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.RecordingDetails
		err = rows.Scan(&r.Id, &r.CallDate, &r.SRC, &r.DST, &r.Duration, &r.BillSec, &r.Disposition, &r.S3FileURL, &r.Office)
		if err != nil {
			return nil, err
		}
		tempURL := strings.Replace(r.S3FileURL,"https://s3.eu-central-1.amazonaws.com/","http://192.168.1.7/",-1)
		tempURL2 := strings.Replace(tempURL,"gsm","wav",-1)
		r.S3FileURL = tempURL2
		results = append(results, r)
	}

	return &results, err
}

func GetRecordingsByNumber(num string) (*[]model.RecordingDetails, error) {
	var results []model.RecordingDetails
	db, err := sql.Open(server, dbURL)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id,calldate,src,dst,duration,billsec,disposition,s3fileurl,office FROM recordings WHERE src LIKE ? OR dst LIKE ?"
	rows, err := db.Query(query,num,num)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.RecordingDetails
		err = rows.Scan(&r.Id, &r.CallDate, &r.SRC, &r.DST, &r.Duration, &r.BillSec, &r.Disposition, &r.S3FileURL, &r.Office)
		if err != nil {
			return nil, err
		}
		tempURL := strings.Replace(r.S3FileURL,"https://s3.eu-central-1.amazonaws.com/","http://192.168.1.7/",-1)
		tempURL2 := strings.Replace(tempURL,"gsm","wav",-1)
		r.S3FileURL = tempURL2
		results = append(results, r)
	}

	return &results,err
}
