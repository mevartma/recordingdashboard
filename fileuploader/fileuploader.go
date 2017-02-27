package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	_ "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	_ "github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
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
var s3recordings []RecordingDetails
var setting ServerConfig

type RecordingDetails struct {
	CallDate       string `json:"calldate"`
	ClId           string `json:"clid"`
	SRC            string `json:"src"`
	DST            string `json:"dst"`
	Duration       string `json:"duration"`
	BillSec        string `json:"billsec"`
	Disposition    string `json:"disposition"`
	AccountCode    string `json:"accountcode"`
	UniqueId       string `json:"uniqueid"`
	DID            string `json:"did"`
	Recording_File string `json:"recordingfile"`
	Disk_File_Path string `json:"disk_file_path"`
	S3_File_URL    string `json:"s_3_file_url"`
	Office         string `json:"office"`
}

type ServerConfig struct {
	Office     string `json:"office"`
	Server_URL string `json:"server_url"`
	AWS_ID     string `json:"aws_id"`
	AWS_Key    string `json:"aws_key"`
}

func init() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatal("Config File Missing. ", err)
		os.Exit(1)
	}
	err = json.Unmarshal(file, &setting)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
		os.Exit(1)
	}
}

func GetAllRecording(date string) (*[]RecordingDetails, error) {
	var results []RecordingDetails
	db, err := sql.Open(server, dbURL)
	if err != nil {
		return &results, err
	}
	defer db.Close()

	query := "SELECT calldate,clid,src,dst,duration,billsec,disposition,accountcode,uniqueid,did,recordingfile FROM cdr WHERE calldate like ?"
	rows, err := db.Query(query, date)
	if err != nil {
		return &results, err
	}

	for rows.Next() {
		var r RecordingDetails
		err = rows.Scan(&r.CallDate, &r.ClId, &r.SRC, &r.DST, &r.Duration, &r.BillSec, &r.Disposition, &r.AccountCode, &r.UniqueId, &r.DID, &r.Recording_File)
		if err != nil {
			return &results, err
		}
		results = append(results, r)
	}

	return &results, err
}

func updateRecords() error {
	now := time.Now()
	newDate := now.Format("2006-01-02")
	newDate += "%"
	rss, err := GetAllRecording(newDate)
	recordings = nil
	for _, rs := range *rss {
		recordings = append(recordings, rs)
	}
	return err
}

func findRecord(recordDate, recordName, officeName string) string {
	var basePath string
	if strings.ToLower(officeName) == "germany" {
		basePath = "/var/spool/asterisk/monitor"
	} else if strings.ToLower(officeName) == "kiev" {
		basePath = "/home/recording"
	}

	s := strings.Split(recordDate, " ")
	date := strings.Split(s[0], "-")
	g := fmt.Sprintf("%s/%s/%s/%s/%s", basePath, date[0], date[1], date[2], recordName)

	return g
}

func Upload2S3() error {
	bucket := aws.String("betamediarecording")
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(setting.AWS_ID, setting.AWS_Key, ""),
		Region:           aws.String("eu-central-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("s3.amazonaws.com"),
	}
	s3Client := s3.New(s3Config)
	/*newSession, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal(err)
		return err
	}*/

	for _, r := range s3recordings {
		file, err := os.Open(r.Disk_File_Path)
		if err != nil {
			log.Println(err)
			return err
		}
		defer file.Close()

		fileInfo, _ := file.Stat()
		var size int64 = fileInfo.Size()

		buffer := make([]byte, size)
		file.Read(buffer)
		fileBytes := bytes.NewReader(buffer)
		fileType := http.DetectContentType(buffer)
		filePath := fmt.Sprintf("/%s/%s", r.Office, r.Recording_File)

		params := &s3.PutObjectInput{
			Bucket:        bucket,
			Key:           aws.String(filePath),
			ACL:           aws.String("public-read"),
			Body:          fileBytes,
			ContentLength: aws.Int64(size),
			ContentType:   aws.String(fileType),
			Metadata: map[string]*string{
				"key": aws.String("MetadataValue"),
			},
		}

		result, err := s3Client.PutObject(params)

		/*uploader := s3manager.NewUploader(newSession)
		_, err = uploader.Upload(&s3manager.UploadInput{
			Bucket: bucket,
			Key:    aws.String(filePath),
			ACL: aws.String("public-read"),
			Body:   file,
			ContentType: aws.String(fileType),
			Metadata: map[string]*string{
				"key": aws.String("MetadataValue"),
			},
		})*/
		if err != nil {
			log.Println(err)
			return err
		}
		fileURL := fmt.Sprintf("https://s3.eu-central-1.amazonaws.com/betamediarecording/%s/%s", r.Office, r.Recording_File)
		r.S3_File_URL = fileURL
		fmt.Println(r.S3_File_URL, awsutil.StringValue(result))
	}

	return nil
}

func UploadToDatabase() error {
	client := http.Client{}
	for _, rec := range s3recordings {
		js, err := json.Marshal(rec)
		if err != nil {
			log.Fatal(err)
			return err
		}
		req, err := http.NewRequest("POST", setting.Server_URL, bytes.NewReader(js))
		if err != nil {
			log.Fatal(err)
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		_, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func main() {
	err := updateRecords()
	if err != nil {
		log.Fatal("Faild to get data from database", err)
		os.Exit(1)
	}

	s3recordings = nil
	for _, record := range recordings {
		if record.Recording_File != "" {
			DiskFilePath := findRecord(record.CallDate, record.Recording_File, setting.Office)
			record.Disk_File_Path = DiskFilePath
			record.Office = setting.Office
			s3recordings = append(s3recordings, record)
		}
	}

	err = Upload2S3()
	if err != nil {
		log.Fatal("Failed to upload files to s3", err)
		os.Exit(1)
	}

	err = UploadToDatabase()
	if err != nil {
		log.Fatal("Failed to upload to database. ", err)
		os.Exit(1)
	}
}
