package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
	"time"
)

type PostRepositories struct {
	db *sql.DB
}

func InitPostRepositories(db *sql.DB) *PostRepositories {
	return &PostRepositories{db}
}

func (r *PostRepositories) CreatePost(post models.Post) (int, error) {
	query := "INSERT INTO post (post_date, post_content, feed_id, user_id) VALUES (?, ?, ?, ?);"

	result, err := r.db.Exec(query,
		time.Now().Format("2006-01-02 15:04:05"),
		post.Content,
		post.FeedID,
		post.UserID,
	)

	if err != nil {
		return -1, fmt.Errorf("erreur lors de la création du post: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("erreur lors de la récupération de l'ID du post: %v", err)
	}

	return int(id), nil
}

func (r *PostRepositories) GetPostsByFeedID(feedID int) ([]models.Post, error) {
	var posts []models.Post
	query := "SELECT post_id, post_date, post_content, feed_id, user_id FROM post WHERE feed_id = ? ORDER BY post_date DESC;"

	rows, err := r.db.Query(query, feedID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des posts: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.Id, &post.Date, &post.Content, &post.FeedID, &post.UserID)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la lecture d'un post: %v", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
} 