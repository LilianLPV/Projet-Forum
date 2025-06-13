package models

type Account struct {
	Id             int
	Username       string
	Password       string
	Email          string
	Role           string
	Profilepicture string
	Bio            string
	Lastconnection string
}
