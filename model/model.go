package model

/*type RecordingDetails struct {
	Id            int64  `json:"id"`
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
	DiskFilePath  string `json:"disk_file_path"`
	S3FileURL     string `json:"s_3_file_url"`
	Office        string `json:"office"`
}*/

type RecordingDetails struct {
	UniqueId      string `json:"unique_id"`
	CallDate      string `json:"call_date"`
	ClId          string `json:"cl_id"`
	SRC           string `json:"src"`
	DST           string `json:"dst"`
	Duration      string `json:"duration"`
	BillSec       string `json:"billsec"`
	Disposition   string `json:"disposition"`
	Cnam          string `json:"cnam"`
	RecordingFile string `json:"recordingfile"`
	FileURL       string `json:"s_3_file_url"`
	ServerURL     string `json:"server_url"`
	Office        string `json:"office"`
}

type RecordingSetting struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

type UserGroups struct {
	GroupName string `json:"groupname"`
}

type UserNameAndPassword struct {
	Username string
	Password string
}

type UserName struct {
	Username string `json:"username"`
}

type UserResult struct {
	Username string `json:"username"`
	Result   bool   `json:"result"`
}

type UserDetails struct {
	Id         int64
	UserName   string
	IpAddress  string
	UserAgent  string
	Cookie     string
	ExpireTime string
}
