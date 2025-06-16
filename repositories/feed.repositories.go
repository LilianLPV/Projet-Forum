package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
	"log"
)

type FeedRepositories struct {
	db *sql.DB
}

func InitFeedRepositories(db *sql.DB) *FeedRepositories {
	return &FeedRepositories{db}
}

func (r *FeedRepositories) CreateFeed(feed models.Feed) (int, error) {
	query := "INSERT INTO feed (feed_title, feed_description, feed_state, feed_creation_date, user_id) VALUES (?, ?, ?, ?, ?);"

	result, err := r.db.Exec(query,
		feed.Title,
		feed.Description,
		feed.State,
		feed.CreatedDate,
		feed.UserID,
	)

	if err != nil {
		log.Printf("Erreur db.Exec lors de la création du feed: %v", err)
		return -1, fmt.Errorf("erreur lors de la création du feed: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Erreur LastInsertId lors de la récupération de l'ID du feed: %v", err)
		return -1, fmt.Errorf("erreur lors de la récupération de l'ID du feed: %v", err)
	}

	log.Printf("Feed créé avec succès, ID: %d", id)
	return int(id), nil
}

func (r *FeedRepositories) ReadById(id int) (models.Feed, error) {
	var feed models.Feed
	query := "SELECT feed_id, feed_title, feed_description, feed_state, feed_creation_date, user_id FROM feed WHERE feed_id = ?;"

	err := r.db.QueryRow(query, id).Scan(
		&feed.Id,
		&feed.Title,
		&feed.Description,
		&feed.State,
		&feed.CreatedDate,
		&feed.UserID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Feed{}, fmt.Errorf("feed non trouvé")
		}
		return models.Feed{}, fmt.Errorf("erreur lors de la lecture du feed: %v", err)
	}

	return feed, nil
}

func (r *FeedRepositories) ReadByUserId(userId int) ([]models.Feed, error) {
	var feeds []models.Feed
	query := "SELECT feed_id, feed_title, feed_description, feed_state, feed_creation_date, user_id FROM feed WHERE user_id = ?;"

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture des feeds: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var feed models.Feed
		err := rows.Scan(
			&feed.Id,
			&feed.Title,
			&feed.Description,
			&feed.State,
			&feed.CreatedDate,
			&feed.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la lecture d'un feed: %v", err)
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}

func (r *FeedRepositories) AddTagToFeed(feedId int, tagId int) error {
	query := "INSERT INTO add_tag (feed_id, tag_id) VALUES (?, ?);"
	_, err := r.db.Exec(query, feedId, tagId)
	if err != nil {
		return fmt.Errorf("erreur lors de l'ajout du tag au feed: %v", err)
	}
	return nil
}

func (r *FeedRepositories) GetFeedTags(feedId int) ([]models.Tag, error) {
	var tags []models.Tag
	query := `
		SELECT t.tag_id, t.tag_name, t.tag_type 
		FROM tag t 
		JOIN add_tag at ON t.tag_id = at.tag_id 
		WHERE at.feed_id = ?;
	`

	rows, err := r.db.Query(query, feedId)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des tags: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.Id, &tag.Name, &tag.Type)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la lecture d'un tag: %v", err)
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *FeedRepositories) UpdateFeedState(feedId int, state string) error {
	query := "UPDATE feed SET feed_state = ? WHERE feed_id = ?;"
	result, err := r.db.Exec(query, state, feedId)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour de l'état du feed: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification des lignes affectées: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("feed non trouvé")
	}

	return nil
}

func (r *FeedRepositories) GetFeedsByTagId(tagId int) ([]models.Feed, error) {
	var feeds []models.Feed
	query := `
		SELECT f.feed_id, f.feed_title, f.feed_description, f.feed_state, f.feed_creation_date, f.user_id 
		FROM feed f 
		JOIN add_tag at ON f.feed_id = at.feed_id 
		WHERE at.tag_id = ?;
	`

	rows, err := r.db.Query(query, tagId)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des feeds: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var feed models.Feed
		err := rows.Scan(
			&feed.Id,
			&feed.Title,
			&feed.Description,
			&feed.State,
			&feed.CreatedDate,
			&feed.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la lecture d'un feed: %v", err)
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
}

func (r *FeedRepositories) ReadAll() ([]models.Feed, error) {
	var feeds []models.Feed
	query := "SELECT feed_id, feed_title, feed_description, feed_state, feed_creation_date, user_id FROM feed ORDER BY feed_creation_date DESC;"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture des feeds: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var feed models.Feed
		err := rows.Scan(
			&feed.Id,
			&feed.Title,
			&feed.Description,
			&feed.State,
			&feed.CreatedDate,
			&feed.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la lecture d'un feed: %v", err)
		}
		feeds = append(feeds, feed)
	}

	return feeds, nil
} 