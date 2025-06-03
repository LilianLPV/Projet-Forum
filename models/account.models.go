package models

import "time"

type Account struct {
	Id             int
	Username       string
	Password       string
	Email          string
	Profilepicture string
	Bio            string
	Lastconnection time.Time
}
