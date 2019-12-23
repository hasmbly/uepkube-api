package controllers

import (
	"time"
	"os"
	"log"
)

func DeleteFile (path string, tm_optional ...int) {
	t := 12

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
