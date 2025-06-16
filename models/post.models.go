package models

type Post struct {
	Id      int    `json:"Id"`
	Date    string `json:"Date"`
	Content string `json:"Content"`
	FeedID  int    `json:"FeedID"`
	UserID  int    `json:"UserID"`
}

type PostWithAuthor struct {
	Post
	AuthorName string `json:"AuthorName"`
}
