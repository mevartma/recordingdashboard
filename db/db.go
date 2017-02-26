package db

import (
	"RecordingDashboard/model"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
)

const (
	dbURL  string = "recordinguser:410QMYbh@tcp(localhost:3306)/recordingdb"
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

func GetAllRecordingsByRange(from, to int64) (*[]model.RecordingDetails, error) {
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
		results = append(results, r)
	}

	return &results, err
}
