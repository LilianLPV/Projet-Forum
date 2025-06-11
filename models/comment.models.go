package models

type Comment struct {
	Id        int
	date      string
	content   string
	commentId int
	postId    int
	userId    int
}
