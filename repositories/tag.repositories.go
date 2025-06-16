package repositories

import (
	"database/sql"
	"fmt"
	"forum/models"
)

type TagRepositories struct {
	db *sql.DB
}

func InitTagRepositories(db *sql.DB) *TagRepositories {
	return &TagRepositories{db}
}

func (r *TagRepositories) CreateTag(tag models.Tag) (int, error) {
	query := "INSERT INTO tag (tag_name, tag_type) VALUES (?, ?);"

	result, err := r.db.Exec(query, tag.Name, tag.Type)
	if err != nil {
		return -1, fmt.Errorf("erreur lors de la création du tag: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("erreur lors de la récupération de l'ID du tag: %v", err)
	}

	return int(id), nil
}

func (r *TagRepositories) ReadById(id int) (models.Tag, error) {
	var tag models.Tag
	query := "SELECT tag_id, tag_name, tag_type FROM tag WHERE tag_id = ?;"

	err := r.db.QueryRow(query, id).Scan(&tag.Id, &tag.Name, &tag.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Tag{}, fmt.Errorf("tag non trouvé")
		}
		return models.Tag{}, fmt.Errorf("erreur lors de la lecture du tag: %v", err)
	}

	return tag, nil
}

func (r *TagRepositories) ReadAll() ([]models.Tag, error) {
	var tags []models.Tag
	query := "SELECT tag_id, tag_name, tag_type FROM tag;"

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture des tags: %v", err)
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

func (r *TagRepositories) GetTagsByFeedId(feedId int) ([]models.Tag, error) {
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

func (r *TagRepositories) GetTagByName(name string) (models.Tag, error) {
	var tag models.Tag
	query := "SELECT tag_id, tag_name, tag_type FROM tag WHERE tag_name = ?;"

	err := r.db.QueryRow(query, name).Scan(&tag.Id, &tag.Name, &tag.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Tag{}, fmt.Errorf("tag non trouvé")
		}
		return models.Tag{}, fmt.Errorf("erreur lors de la lecture du tag: %v", err)
	}

	return tag, nil
} 