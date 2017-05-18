package db

import (
	"RecordingDashboard/model"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"fmt"
	"github.com/aws/aws-sdk-go/private/protocol/query"
)

const (
	db_local_URL  string = "root:@tcp(localhost:3306)/recordingdb"
	db_germany_URL  string = "recording:recording@tcp(192.168.50.14:3306)/asteriskcdrdb"
	db_kiev_URL  string = "recording:recording@tcp(localhost:3306)/asteriskcdrdb"
	server string = "mysql"
)

var db *sql.DB
var err error

/*func UpdateRecording(r model.RecordingDetails, action string) error {
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
}*/

func GetSessionId(s string) (bool, *model.UserDetails, error) {
	err = nil
	var results model.UserDetails
	db, err = sql.Open(server, db_local_URL)
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

func GetAllSessions() (*[]model.UserDetails, error) {
	err = nil
	var result []model.UserDetails
	db, err = sql.Open(server, db_local_URL)
	defer db.Close()
	query := "SELECT * FROM userssessions"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.UserDetails
		err = rows.Scan(&r.Id, &r.UserName, &r.IpAddress, &r.UserAgent, &r.Cookie, &r.ExpireTime)
		if err != nil {
			return nil, err
		}
		result = append(result,r)
	}

	return &result,err
}

func DeleteSessionId(s string) error {
	err = nil
	db, err = sql.Open(server, db_local_URL)
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM userssessions WHERE cookie = ?")
	_, err = stmt.Exec(s)

	return err
}

func UpdateUser(r model.UserDetails, action string) error {
	err = nil
	db, err = sql.Open(server, db_local_URL)
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

func GetAllRecordings(office string) (*[]model.RecordingDetails, error) {
	var results []model.RecordingDetails
	err = nil
	switch office {
	case "germany":
		db, err = sql.Open(server, db_germany_URL)
	case "'kiev":
		db, err = sql.Open(server, db_kiev_URL)
	}
	//db, err := sql.Open(server, dbURL)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT calldate,clid,src,dst,duration,billsec,disposition,s3fileurl,office FROM recordings"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.RecordingDetails
		err = rows.Scan(&r.CallDate, &r.ClId, &r.SRC, &r.DST, &r.Duration, &r.BillSec, &r.Disposition, &r.S3FileURL, &r.Office)
		if err != nil {
			return nil, err
		}
		if r.Disposition != "NO ANSWER" {
			tempURL := strings.Replace(r.S3FileURL, "https://s3.eu-central-1.amazonaws.com/", "http://192.168.1.7/", -1)
			tempURL2 := strings.Replace(tempURL, "gsm", "mp3", -1)
			r.S3FileURL = tempURL2
			results = append(results, r)
		}
	}

	return &results, err
}

/*func GetRecordingsByRange(from, to int64, office string) (*[]model.RecordingDetails, error) {
	var results []model.RecordingDetails
	err = nil
	switch office {
	case "germany":
		db, err = sql.Open(server, db_germany_URL)
	case "'kiev":
		db, err = sql.Open(server, db_kiev_URL)
	}
	//db, err := sql.Open(server, dbURL)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT calldate,clid,src,dst,duration,billsec,disposition,s3fileurl,office FROM cdr WHERE id BETWEEN ? AND ?"
	rows, err := db.Query(query, from, to)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.RecordingDetails
		err = rows.Scan(&r.CallDate, &r.ClId, &r.SRC, &r.DST, &r.Duration, &r.BillSec, &r.Disposition, &r.S3FileURL, &r.Office)
		if err != nil {
			return nil, err
		}
		if r.Disposition != "NO ANSWER" {
			tempURL := strings.Replace(r.S3FileURL, "https://s3.eu-central-1.amazonaws.com/", "http://192.168.1.7/", -1)
			tempURL2 := strings.Replace(tempURL, "gsm", "mp3", -1)
			r.S3FileURL = tempURL2
			results = append(results, r)
		}
	}

	return &results, err
}*/

/*func GetRecordingsByNumber(num string, office string) (*[]model.RecordingDetails, error) {
	var results []model.RecordingDetails
	err = nil
	switch office {
	case "germany":
		db, err = sql.Open(server, db_germany_URL)
	case "'kiev":
		db, err = sql.Open(server, db_kiev_URL)
	}
	//db, err := sql.Open(server, dbURL)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT calldate,clid,src,dst,duration,billsec,disposition,s3fileurl,office FROM cdr WHERE src LIKE ? OR dst LIKE ?"
	rows, err := db.Query(query,num,num)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var r model.RecordingDetails
		err = rows.Scan(&r.CallDate, &r.ClId, &r.SRC, &r.DST, &r.Duration, &r.BillSec, &r.Disposition, &r.FileURL, &r.Office)
		fmt.Println(r)
		if err != nil {
			return nil, err
		}
		if r.Disposition != "NO ANSWER" {
			tempURL := strings.Replace(r.S3FileURL, "https://s3.eu-central-1.amazonaws.com/", "http://192.168.1.7/", -1)
			tempURL2 := strings.Replace(tempURL, "gsm", "mp3", -1)
			r.S3FileURL = tempURL2
			results = append(results, r)
		}
	}

	return &results,err
}*/

func GetRecording(num, date1, date2, office string) (*[]model.RecordingDetails, error) {
	var results []model.RecordingDetails
	err = nil

	switch office {
	case "germany":
		db, err = sql.Open(server, db_germany_URL)
	case "'kiev":
		db, err = sql.Open(server, db_kiev_URL)
	default:
		err = error("Choose Right Office")
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT calldate,clid,src,dst,duration,billsec,disposition,cnam,recordingfile WHERE src LIKE ? OR dst LIKE ? AND disposition like 'ANSWERED' AND calldate BETWEEN ? AND ?"
	rows, err := db.Query(query,num,num,date1,date2)
	for rows.Next() {
		var r model.RecordingDetails
		err = rows.Scan(&r.CallDate,&r.ClId,&r.SRC,&r.DST,&r.Duration,&r.BillSec,&r.Disposition,&r.Cnam,&r.RecordingFile,&r.UniqueId)
		r.Office = office
		switch office {
		case "gemany":
			r.FileURL = fmt.Sprintf("https://s3.eu-central-1.amazonaws.com/betamediarecording/Germany/%s", r.RecordingFile)
			r.ServerURL = fmt.Sprintf("http://192.168.1.7/betamediarecording/Germany/%s", r.RecordingFile)
		case "kiev":
			r.FileURL = fmt.Sprintf("https://s3.eu-central-1.amazonaws.com/betamediarecording/Kiev/%s", r.RecordingFile)
			r.ServerURL = fmt.Sprintf("http://192.168.1.7/betamediarecording/Kiev/%s", r.RecordingFile)
		}

		results = append(results,r)
	}

	return &results,err
}
