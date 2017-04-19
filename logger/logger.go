package logger

/*import (
	"log"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

const capacity = 32768

type Worker struct {
	FileRoot string
	Buffer []byte
	Position int
}

func NewWorker(id int) (w *Worker) {
	return &Worker{
		FileRoot: "/var/log/recordingsystem/" + strconv.Itoa(id) + "_",
		Buffer: make([]byte, capacity),
	}
}

func (w *Worker) Work (channel chan []byte) {
	for {
		event := <- channel
		length := len(event)

		if length > capacity {
			log.Println("To big")
			continue
		}
		if(length + w.Position) > capacity {
			w.Save()
		}
		copy(w.Buffer[w.Position:],event)
		w.Position += length
	}
}

func (w *Worker) Save() {
	if w.Position == 0 { return }
	f, _ := ioutil.TempFile("", "logs_")
	f.Write(w.Buffer[0:w.Position])
	f.Close()
	os.Rename(f.Name(),w.FileRoot + strconv.FormatInt(time.Now().UnixNano(), 10))
	w.Position = 0
}*/
