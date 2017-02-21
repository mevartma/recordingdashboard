package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	dbURL  string = "root:@tcp(localhost:3306)/asteriskcdrdb"
	server string = "mysql"
)

var recordings []RecordingDetails
var s3recordings []S3RecordingFileDetails
var office = string("Germany")
var postURL = string("")
var count int64

type RecordingDetails struct {
	CallDate      string `json:"calldate"`
	ClId          string `json:"clid"`
	SRC           string `json:"src"`
	DST           string `json:"dst"`
	Duration      string `json:"duration"`
	BillSec       string `json:"billsec"`
	Disposition   string `json:"disposition"`
	AccountCode   string `json:"accountcode"`
	UniqueId      string `json:"uniqueid"`
	DID           string `json:"did"`
	RecordingFile string `json:"recordingfile"`
}

type S3RecordingFileDetails struct {
	RecordingDetails
	DiskFilePath string `json:"disk_file_path"`
	S3FileURL    string `json:"s_3_file_url"`
	Office       string `json:"office"`
}

func GetAllRecodring(date string) (*[]RecordingDetails, error) {
	var results []RecordingDetails
	db, err := sql.Open(server, dbURL)
	if err != nil {
		return &results, err
	}
	defer db.Close()

	query := "SELECT calldate,clid,src,dst,duration,billsec,disposition,accountcode,uniqueid,did,recordingfile FROM cdr WHERE calldate like '?%'"
	rows, err := db.Query(query,date)
	if err != nil {
		return &results, err
	}

	for rows.Next() {
		var r RecordingDetails
		err = rows.Scan(&r.CallDate, &r.ClId, &r.SRC, &r.DST, &r.Duration, &r.BillSec, &r.Disposition, &r.AccountCode, &r.UniqueId, &r.DID, &r.RecordingFile)
		if err != nil {
			return &results, err
		}
		results = append(results, r)
	}

	return &results, err
}

func updateRecords() error {
	now := time.Now()
	newDate := fmt.Sprintf("%s",now.Format("2006-01-02"))
	rss, err := GetAllRecodring(newDate)
	recordings = nil
	for _, rs := range *rss {
		recordings = append(recordings, rs)
	}
	return err
}

func findRecord(recoredDate, recordName, officeName string) string {
	var basePath string
	if officeName == "Germany" {
		basePath = "/var/spool/asterisk/monitor"
	} else if officeName == "Kiev" {
		basePath = "/home/recording"
	}

	s := strings.Split(recoredDate, " ")
	date := strings.Split(s[0], "-")
	g := fmt.Sprintf("%s/%s/%s/%s/%s", basePath, date[0], date[1], date[2], recordName)

	return g
}

func Upload2S3() error {
	bucket := aws.String("betamediarecording")
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("AKIAJBORGMMA5NY2GRFQ", "8YcCRVf6h+DG2qCVLFM6y1PUn8VfJhr8R1oelrBV", ""),
		Region:           aws.String("eu-central-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, r := range s3recordings {
		count++
		if count == 30 {
			return nil
		} else {
			file, err := os.Open(r.DiskFilePath)
			if err != nil {
				log.Println(err)
				return err
			}
			defer file.Close()

			uploader := s3manager.NewUploader(newSession)
			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: bucket,
				Key:    aws.String(r.RecordingFile),
				Body:   file,
			})
			if err != nil {
				log.Println(err)
				return err
			}
			fileURL := fmt.Sprintf("https://s3.eu-central-1.amazonaws.com/betamediarecording/%s/%s", r.Office, r.RecordingFile)
			r.S3FileURL = fileURL
			fmt.Println(r.S3FileURL)
			time.Sleep(4 * time.Second)
		}
	}

	return nil
}

func UploadToDatabase() (err error) {
	client := http.Client{}
	for _, rec := range s3recordings {
		js, err := json.Marshal(rec)
		if err != nil {
			log.Fatal(err)
			return
		}
		req, err := http.NewRequest("POST", postURL, bytes.NewReader(js))
		if err != nil {
			log.Fatal(err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		_, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	return nil
}

func main() {
	err := updateRecords()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	s3recordings = nil
	for _, record := range recordings {
		if record.RecordingFile != "" {
			diskfilepath := findRecord(record.CallDate, record.RecordingFile, office)
			var s3r= S3RecordingFileDetails{record, diskfilepath, "", office}
			s3recordings = append(s3recordings, s3r)
		}
	}

	err = Upload2S3()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
