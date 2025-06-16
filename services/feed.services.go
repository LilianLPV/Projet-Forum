package services

import (
	"database/sql"
	"fmt"
	"forum/models"
	"forum/repositories"
	"strings"
	"time"
)

type FeedService struct {
	FeedRepo    *repositories.FeedRepositories
	TagRepo     *repositories.TagRepositories
	PostRepo    *repositories.PostRepositories
	AccountRepo *repositories.AccountRepositories
}

func InitFeedService(db *sql.DB) *FeedService {
	return &FeedService{
		FeedRepo:    repositories.InitFeedRepositories(db),
		TagRepo:     repositories.InitTagRepositories(db),
		PostRepo:    repositories.InitPostRepositories(db),
		AccountRepo: repositories.InitAccountRepositories(db),
	}
}

func (s *FeedService) CreateFeed(feed models.Feed, tags []string) (int, error) {
	if feed.Title == "" || feed.Description == "" {
		return -1, fmt.Errorf("titre et description requis")
	}

	feed.CreatedDate = time.Now().Format("2006-01-02 15:04:05")
	if feed.State == "" {
		feed.State = "open" // Par défaut, le feed est ouvert
	}

	// Vérifier que l'état est valide
	if feed.State != "open" && feed.State != "closed" {
		return -1, fmt.Errorf("état invalide, doit être 'open' ou 'closed'")
	}

	feedId, err := s.FeedRepo.CreateFeed(feed)
	if err != nil {
		return -1, fmt.Errorf("erreur lors de la création du feed: %v", err)
	}

	// Ajouter les tags
	if len(tags) > 0 {
		for _, tagName := range tags {
			tagName = strings.TrimSpace(tagName)
			if tagName == "" {
				continue
			}

			// Vérifier si le tag existe déjà
			tag, err := s.TagRepo.GetTagByName(tagName)
			if err != nil {
				// Créer le tag s'il n'existe pas
				tag = models.Tag{
					Name: tagName,
					Type: "user", // Type par défaut pour les tags créés par les utilisateurs
				}
				tagId, err := s.TagRepo.CreateTag(tag)
				if err != nil {
					return -1, fmt.Errorf("erreur lors de la création du tag: %v", err)
				}
				tag.Id = tagId
			}

			// Associer le tag au feed
			err = s.FeedRepo.AddTagToFeed(feedId, tag.Id)
			if err != nil {
				return -1, fmt.Errorf("erreur lors de l'ajout du tag au feed: %v", err)
			}
		}
	}

	return feedId, nil
}

func (s *FeedService) GetFeedById(id int) (models.Feed, []models.Tag, error) {
	if id <= 0 {
		return models.Feed{}, nil, fmt.Errorf("ID de feed invalide")
	}

	feed, err := s.FeedRepo.ReadById(id)
	if err != nil {
		return models.Feed{}, nil, fmt.Errorf("erreur lors de la récupération du feed: %v", err)
	}

	tags, err := s.FeedRepo.GetFeedTags(id)
	if err != nil {
		return feed, nil, fmt.Errorf("erreur lors de la récupération des tags: %v", err)
	}

	return feed, tags, nil
}

func (s *FeedService) GetFeedsByUserId(userId int) ([]models.Feed, error) {
	if userId <= 0 {
		return nil, fmt.Errorf("ID utilisateur invalide")
	}

	feeds, err := s.FeedRepo.ReadByUserId(userId)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des feeds: %v", err)
	}

	return feeds, nil
}

func (s *FeedService) UpdateFeedState(feedId int, state string) error {
	if feedId <= 0 {
		return fmt.Errorf("ID de feed invalide")
	}

	if state != "open" && state != "closed" {
		return fmt.Errorf("état invalide, doit être 'open' ou 'closed'")
	}

	err := s.FeedRepo.UpdateFeedState(feedId, state)
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour de l'état du feed: %v", err)
	}

	return nil
}

func (s *FeedService) GetAllFeeds() ([]models.FeedWithAuthor, error) {
	feeds, err := s.FeedRepo.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération de tous les feeds: %v", err)
	}

	var feedsWithAuthor []models.FeedWithAuthor
	for _, feed := range feeds {
		authorName, err := s.GetUsernameByUserID(feed.UserID)
		if err != nil {
			// Log l'erreur mais continue sans nom d'auteur si l'utilisateur n'est pas trouvé
			authorName = "Utilisateur inconnu"
		}
		feedsWithAuthor = append(feedsWithAuthor, models.FeedWithAuthor{
			Feed:       feed,
			AuthorName: authorName,
		})
	}
	return feedsWithAuthor, nil
}

func (s *FeedService) GetFeedTags(feedId int) ([]models.Tag, error) {
	tags, err := s.FeedRepo.GetFeedTags(feedId)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des tags: %v", err)
	}
	return tags, nil
}

func (s *FeedService) GetTagByName(name string) (models.Tag, error) {
	return s.TagRepo.GetTagByName(name)
}

func (s *FeedService) GetFeedsByTagId(tagId int) ([]models.Feed, error) {
	return s.FeedRepo.GetFeedsByTagId(tagId)
}

func (s *FeedService) CreatePost(post models.Post) error {
	if post.Content == "" {
		return fmt.Errorf("le contenu du post ne peut pas être vide")
	}
	if post.FeedID <= 0 {
		return fmt.Errorf("ID de feed invalide pour le post")
	}
	if post.UserID <= 0 {
		return fmt.Errorf("ID utilisateur invalide pour le post")
	}

	_, err := s.PostRepo.CreatePost(post)
	if err != nil {
		return fmt.Errorf("erreur lors de la création du post: %v", err)
	}

	return nil
}

func (s *FeedService) GetFeedPosts(feedID int) ([]models.PostWithAuthor, error) {
	if feedID <= 0 {
		return nil, fmt.Errorf("ID de feed invalide")
	}
	posts, err := s.PostRepo.GetPostsByFeedID(feedID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des posts du feed: %v", err)
	}

	var postsWithAuthor []models.PostWithAuthor
	for _, post := range posts {
		authorName, err := s.GetUsernameByUserID(post.UserID)
		if err != nil {
			authorName = "Utilisateur inconnu"
		}
		postsWithAuthor = append(postsWithAuthor, models.PostWithAuthor{
			Post:       post,
			AuthorName: authorName,
		})
	}
	return postsWithAuthor, nil
}

// Fonction utilitaire pour récupérer le nom d'utilisateur
func (s *FeedService) GetUsernameByUserID(userID int) (string, error) {
	account, err := s.AccountRepo.ReadById(userID)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la récupération de l'utilisateur: %v", err)
	}
	return account.Username, nil
}

func (s *FeedService) GetFeedByIdWithAuthor(id int) (models.FeedWithAuthor, []models.Tag, error) {
	feed, tags, err := s.GetFeedById(id)
	if err != nil {
		return models.FeedWithAuthor{}, nil, err
	}
	authorName, err := s.GetUsernameByUserID(feed.UserID)
	if err != nil {
		authorName = "Utilisateur inconnu"
	}
	return models.FeedWithAuthor{Feed: feed, AuthorName: authorName}, tags, nil
}

func (s *FeedService) GetFeedsByUserIdWithAuthor(userID int) ([]models.FeedWithAuthor, error) {
	feeds, err := s.FeedRepo.ReadByUserId(userID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des feeds par utilisateur: %v", err)
	}

	var feedsWithAuthor []models.FeedWithAuthor
	for _, feed := range feeds {
		authorName, err := s.GetUsernameByUserID(feed.UserID)
		if err != nil {
			authorName = "Utilisateur inconnu"
		}
		feedsWithAuthor = append(feedsWithAuthor, models.FeedWithAuthor{
			Feed:       feed,
			AuthorName: authorName,
		})
	}
	return feedsWithAuthor, nil
}

func (s *FeedService) GetFeedsByTagIdWithAuthor(tagID int) ([]models.FeedWithAuthor, error) {
	feeds, err := s.FeedRepo.GetFeedsByTagId(tagID)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération des feeds par tag: %v", err)
	}

	var feedsWithAuthor []models.FeedWithAuthor
	for _, feed := range feeds {
		authorName, err := s.GetUsernameByUserID(feed.UserID)
		if err != nil {
			authorName = "Utilisateur inconnu"
		}
		feedsWithAuthor = append(feedsWithAuthor, models.FeedWithAuthor{
			Feed:       feed,
			AuthorName: authorName,
		})
	}
	return feedsWithAuthor, nil
} 