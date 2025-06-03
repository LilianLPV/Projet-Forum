package models

import "time"

type post struct {
	Id      int
	date    time.Time
	content string
	feedid  int
	userid  int
}
