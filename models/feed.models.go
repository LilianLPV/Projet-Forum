package models

import "time"

type Feed struct {
	Id          int
	Title       string
	Description string
	State       string
	CreatedDate time.Time
	userid      int
}
