package models

import (
	"time"
)

// OK
type Jn struct {
	Err 	bool 	 	`json:"error" example:"false"`
	Msg 	interface{} `json:"message"`
}

type Jn1 struct {
	Err 	bool 	 	`json:"error" example:"false"`
	Msg 	interface{} `json:"message"`
	Id 		int 		`json:"id"`
}

// NOT OK
type HTTPError struct {
        Code    	int   		`json:"code"`
        Message  	string		`json:"message"`
        Times 		time.Time	`json:"timestamp"`
}