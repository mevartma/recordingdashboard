package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

//------------------------------------------Workers------------------------------------------------//

var (
	WorkQueue   = make(chan TypeOfWork, 200)
	WorkerQueue chan chan TypeOfWork
	client      = &http.Client{}
	bucket      = aws.String("betamediarecording")
	s3Client    *s3.S3
)

type Worker struct {
	Id          int
	Work        chan TypeOfWork
	WorkerQueue chan chan TypeOfWork
	QuitChan    chan bool
}

type TypeOfWork struct {
	RecordingDetails
	WorkType string
}

func NewWorker(id int, workerQueue chan chan TypeOfWork) Worker {
	worker := Worker{
		Id:          id,
		Work:        make(chan TypeOfWork),
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
				if work.WorkType == "s3" {
					if work.Recording_File != "" {
						DiskFilePath := findRecord(work.CallDate, work.Recording_File, setting.Office)
						work.Disk_File_Path = DiskFilePath
						work.Office = setting.Office

						file, err := os.Open(work.Disk_File_Path)
						if err != nil {
							log.Println(err)
						}

						fileInfo, _ := file.Stat()
						var size int64 = fileInfo.Size()
						buffer := make([]byte, size)
						file.Read(buffer)
						fileBytes := bytes.NewReader(buffer)
						fileType := http.DetectContentType(buffer)
						filePath := fmt.Sprintf("/%s/%s", work.Office, work.Recording_File)
						file.Close()

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
						if err != nil {
							log.Println(err)
						}

						fileURL := fmt.Sprintf("https://s3.eu-central-1.amazonaws.com/betamediarecording/%s/%s", work.Office, work.Recording_File)
						work.S3_File_URL = fileURL
						r := RecordingDetails{
							CallDate:       work.CallDate,
							ClId:           work.ClId,
							SRC:            work.SRC,
							DST:            work.DST,
							Duration:       work.Duration,
							BillSec:        work.BillSec,
							Disposition:    work.Disposition,
							AccountCode:    work.AccountCode,
							UniqueId:       work.UniqueId,
							DID:            work.DID,
							Recording_File: work.Recording_File,
							S3_File_URL:    work.S3_File_URL,
							Office:         work.Office,
						}
						s3ProdRecording = append(s3ProdRecording, r)
						fmt.Println(work.S3_File_URL, awsutil.StringValue(result.ETag))
					}
				} else if work.WorkType == "db" {
					r := RecordingDetails{
						CallDate:       work.CallDate,
						ClId:           work.ClId,
						SRC:            work.SRC,
						DST:            work.DST,
						Duration:       work.Duration,
						BillSec:        work.BillSec,
						Disposition:    work.Disposition,
						AccountCode:    work.AccountCode,
						UniqueId:       work.UniqueId,
						DID:            work.DID,
						Recording_File: work.Recording_File,
						S3_File_URL:    work.S3_File_URL,
						Office:         work.Office,
					}
					js, err := json.Marshal(r)
					if err != nil {
						log.Fatal(err)
					}

					req, err := http.NewRequest("POST", setting.Server_URL, bytes.NewReader(js))
					if err != nil {
						log.Fatal(err)
					}

					req.Header.Set("Content-Type", "application/json")
					_, err = client.Do(req)
					if err != nil {
						log.Fatal(err)
					}
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
	WorkerQueue = make(chan chan TypeOfWork, nWorkers)

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

var s3ProdRecording []RecordingDetails
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

	startDispatcher(4)

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(setting.AWS_ID, setting.AWS_Key, ""),
		Region:           aws.String("eu-central-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("s3.eu-central-1.amazonaws.com"),
	}

	s3Session, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal(err)
	}

	s3Client = s3.New(s3Session)
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
	now := time.Now().AddDate(0, 0, -1)
	newDate := now.Format("2006-01-02")
	newDate += "%"
	rss, err := GetAllRecording(newDate)
	s3recordings = nil
	for _, rs := range *rss {
		s3recordings = append(s3recordings, rs)
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

func main() {
	err := updateRecords()
	if err != nil {
		log.Fatal("Faild to get data from database", err)
		os.Exit(1)
	}

	/*bucket := aws.String("betamediarecording")
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(setting.AWS_ID, setting.AWS_Key, ""),
		Region:           aws.String("eu-central-1"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String("s3.eu-central-1.amazonaws.com"),
	}

	s3Session, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal(err)
	}

	s3Client := s3.New(s3Session)*/

	for _, record := range s3recordings {
		r := TypeOfWork{record, "s3"}
		WorkQueue <- r
		/*if record.Recording_File != "" {
			DiskFilePath := findRecord(record.CallDate, record.Recording_File, setting.Office)
			record.Disk_File_Path = DiskFilePath
			record.Office = setting.Office

			file, err := os.Open(record.Disk_File_Path)
			if err != nil {
				log.Println(err)
			}

			fileInfo, _ := file.Stat()
			var size int64 = fileInfo.Size()
			buffer := make([]byte, size)
			file.Read(buffer)
			fileBytes := bytes.NewReader(buffer)
			fileType := http.DetectContentType(buffer)
			filePath := fmt.Sprintf("/%s/%s", record.Office, record.Recording_File)
			file.Close()

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
			if err != nil {
				log.Println(err)
			}

			fileURL := fmt.Sprintf("https://s3.eu-central-1.amazonaws.com/betamediarecording/%s/%s", record.Office, record.Recording_File)
			record.S3_File_URL = fileURL
			s3ProdRecording = append(s3ProdRecording, record)
			fmt.Println(record.S3_File_URL, awsutil.StringValue(result.ETag))
		}*/
	}

	for _, record := range s3ProdRecording {
		r := TypeOfWork{record, "db"}
		WorkQueue <- r
		/*js, err := json.Marshal(record)
		if err != nil {
			log.Fatal(err)
		}
		req, err := http.NewRequest("POST", setting.Server_URL, bytes.NewReader(js))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		_, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}*/
	}

	os.Exit(1)
}
