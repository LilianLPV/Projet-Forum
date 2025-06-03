package models

import "time"

type Comment struct {
	Id        int
	date      time.Time
	content   string
	commentId int
	postId    int
	userId    int
}
