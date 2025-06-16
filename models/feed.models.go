package models

type Feed struct {
	Id          int
	Title       string
	Description string
	State       string
	CreatedDate string
	UserID      int
}

type FeedWithAuthor struct {
	Feed
	AuthorName string
}
