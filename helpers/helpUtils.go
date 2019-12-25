package helpers

import (
	"time"
	"os"
	"log"
	"encoding/base64"	
	"encoding/json"	
)

func DeleteFile (path string, tm_optional ...int) {
	t := 1

	if len(tm_optional) > 0 { t = tm_optional[0] }
	
	// duration
	duration := time.Duration(t)*time.Hour
	// duration := time.Duration(t)*time.Second

	// set & start timers
	timer1 := time.NewTimer(duration)
	<-timer1.C

	// after timers end. delete file
	if err := os.Remove(path); err != nil {
	    log.Println("error deleting file: %v", err)
	}	

	log.Println("file pdf deleted..")

}

func CreateFile (tmpPath string, blob string) bool {
	if blob != "" {

		dec, err := base64.StdEncoding.DecodeString(blob)
		if err != nil {
		    log.Println("error decoding: %v", err)
		}

		// open | create
		f, err := os.OpenFile(tmpPath, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
		    log.Println("error opening file: %v", err)
		}
		defer f.Close()

		// write
		if _, err := f.Write(dec); err != nil {
		    log.Println("error creating file: %v", err)
		}
		if err := f.Sync(); err != nil {
		    log.Println("error Sync file: %v", err)
		}

		go DeleteFile(tmpPath) // <- put to the channel and go routines/thread	
		// go to begginng of file
		// f.Seek(0, 0)
		// output file contents
		// io.Copy(os.Stdout, f)
		
		// delete file with interval time
		// default interval = 12 Hours
	
		return true

	} else {	
		
		return false
	}

	return false
}

func FetchPost (i interface{}) {
	b, err := json.MarshalIndent(i, "", " ")
	if err != nil {
	    log.Println("error:", err)
	}
	os.Stdout.Write(b)
}