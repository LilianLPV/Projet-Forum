package controllers

import (
	"encoding/json"
	"forum/auth"
	"forum/models"
	"forum/services"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type FeedController struct {
	service  *services.FeedService
	template *template.Template
}

func InitFeedController(service *services.FeedService, template *template.Template) *FeedController {
	return &FeedController{service: service, template: template}
}

func (c *FeedController) FeedRouter(r *mux.Router) {
	r.HandleFunc("/feeds/create", c.CreateFeedForm).Methods("GET")
	r.HandleFunc("/feeds/all", c.GetAllFeeds).Methods("GET")
	r.HandleFunc("/feeds", c.CreateFeed).Methods("POST")
	r.HandleFunc("/feeds/{id}", c.GetFeed).Methods("GET")
	r.HandleFunc("/feeds/user", c.GetUserFeeds).Methods("GET")
	r.HandleFunc("/feeds/{id}/state", c.UpdateFeedState).Methods("PUT")
	r.HandleFunc("/feeds/tag/{tag}", c.GetFeedsByTag).Methods("GET")
	r.HandleFunc("/posts/create", c.CreatePost).Methods("POST")
	r.HandleFunc("/posts/feed/{id}", c.GetPostsByFeed).Methods("GET")
}

func (c *FeedController) CreateFeedForm(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil || token.Value == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	_, err = auth.ValidateToken(token.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	c.template.ExecuteTemplate(w, "create_feed", nil)
}

func (c *FeedController) CreateFeed(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil || token.Value == "" {
		log.Printf("Erreur CreateFeed: Token manquant ou invalide: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Non autorisé: " + err.Error()})
		return
	}

	claims, err := auth.ValidateToken(token.Value)
	if err != nil {
		log.Printf("Erreur CreateFeed: Validation du token échouée: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Token invalide: " + err.Error()})
		return
	}

	// Vérifier que les champs requis sont présents
	title := r.FormValue("title")
	description := r.FormValue("description")
	state := r.FormValue("state")

	if title == "" || description == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Titre et description requis"})
		return
	}

	feed := models.Feed{
		Title:       title,
		Description: description,
		State:       state,
		UserID:      claims.Sub,
	}

	// Récupérer les tags
	tagsStr := r.FormValue("tags")
	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
	}

	feedId, err := c.service.CreateFeed(feed, tags)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	if feedId <= 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "La création du feed a échoué: ID invalide"})
		return
	}

	response := struct {
		FeedId int `json:"feed_id"`
	}{
		FeedId: feedId,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (c *FeedController) GetFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	feedIdStr := vars["id"]
	feedId, err := strconv.Atoi(feedIdStr)
	if err != nil {
		log.Printf("Erreur GetFeed: ID de feed invalide: %v", err)
		http.Error(w, "ID de feed invalide", http.StatusBadRequest)
		return
	}

	feedWithAuthor, tags, err := c.service.GetFeedByIdWithAuthor(feedId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	posts, err := c.service.GetFeedPosts(feedId)
	if err != nil {
		log.Printf("Erreur GetFeed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUserID := 0 // Default to 0 (unauthenticated)
	token, err := r.Cookie("token")
	if err == nil && token.Value != "" {
		claims, err := auth.ValidateToken(token.Value)
		if err == nil {
			currentUserID = claims.Sub
		}
	}

	data := struct {
		Feed          models.FeedWithAuthor
		Tags          []models.Tag
		Posts         []models.PostWithAuthor
		CurrentUserID int
	}{
		Feed:          feedWithAuthor,
		Tags:          tags,
		Posts:         posts,
		CurrentUserID: currentUserID,
	}

	c.template.ExecuteTemplate(w, "feed", data)
}

func (c *FeedController) GetUserFeeds(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil || token.Value == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(token.Value)
	if err != nil {
		http.Error(w, "Token invalide", http.StatusUnauthorized)
		return
	}

	feeds, err := c.service.GetFeedsByUserIdWithAuthor(claims.Sub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.template.ExecuteTemplate(w, "user_feeds", struct {
		Feeds []models.FeedWithAuthor
	}{
		Feeds: feeds,
	})
}

func (c *FeedController) UpdateFeedState(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil || token.Value == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	_, err = auth.ValidateToken(token.Value)
	if err != nil {
		http.Error(w, "Token invalide", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	feedIdStr := vars["id"]
	feedId, err := strconv.Atoi(feedIdStr)
	if err != nil {
		http.Error(w, "ID de feed invalide", http.StatusBadRequest)
		return
	}

	var request struct {
		State string `json:"state"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Format de requête invalide", http.StatusBadRequest)
		return
	}

	if request.State != "open" && request.State != "closed" {
		http.Error(w, "État invalide", http.StatusBadRequest)
		return
	}

	err = c.service.UpdateFeedState(feedId, request.State)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *FeedController) GetFeedsByTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagName := vars["tag"]

	// Récupérer le tag par son nom
	tag, err := c.service.GetTagByName(tagName)
	if err != nil {
		http.Error(w, "Tag non trouvé", http.StatusNotFound)
		return
	}

	// Récupérer les feeds associés au tag
	feeds, err := c.service.GetFeedsByTagIdWithAuthor(tag.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.template.ExecuteTemplate(w, "feeds_by_tag", struct {
		TagName string
		Feeds   []models.FeedWithAuthor
	}{
		TagName: tagName,
		Feeds:   feeds,
	})
}

func (c *FeedController) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := c.service.GetAllFeeds()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.template.ExecuteTemplate(w, "all_feeds", struct {
		Feeds []models.FeedWithAuthor
	}{
		Feeds: feeds,
	})
}

func (c *FeedController) CreatePost(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("token")
	if err != nil || token.Value == "" {
		http.Error(w, "Non autorisé", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateToken(token.Value)
	if err != nil {
		http.Error(w, "Token invalide", http.StatusUnauthorized)
		return
	}

	feedIdStr := r.FormValue("feedId")
	content := r.FormValue("content")

	feedId, err := strconv.Atoi(feedIdStr)
	if err != nil {
		http.Error(w, "ID de feed invalide", http.StatusBadRequest)
		return
	}

	if content == "" {
		http.Error(w, "Le contenu du post ne peut pas être vide", http.StatusBadRequest)
		return
	}

	post := models.Post{
		Content: content,
		FeedID:  feedId,
		UserID:  claims.Sub,
	}

	err = c.service.CreatePost(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Post créé avec succès"})
}

func (c *FeedController) GetPostsByFeed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	feedIdStr := vars["id"]
	feedId, err := strconv.Atoi(feedIdStr)
	if err != nil {
		log.Printf("Erreur GetPostsByFeed: ID de feed invalide: %v", err)
		http.Error(w, "ID de feed invalide", http.StatusBadRequest)
		return
	}

	posts, err := c.service.GetFeedPosts(feedId)
	if err != nil {
		log.Printf("Erreur GetPostsByFeed: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Posts récupérés pour feed ID %d: %+v", feedId, posts)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
} 