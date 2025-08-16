package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func PostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getPost(w, r)
	case http.MethodPost:
		createPost(w, r)
	case http.MethodPut:
		updatePost(w, r)
	case http.MethodDelete:
		deletePost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getPost(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("id")
	if postIdStr == "" {
		getPosts(w, r)

		return
	}

	postId, err := strconv.Atoi(postIdStr)
	if err != nil || postId <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	query := `SELECT * FROM posts WHERE id = $1`
	var post Post
	if err := DB.QueryRow(context.Background(), query, postId).Scan(&post.ID, &post.Title, &post.Body, &post.UserId); err != nil {
		http.Error(w, "Failed to get post from DB", http.StatusInternalServerError)
		return
	}

	CustomJsonResponse(w, http.StatusOK, post)
}

func getPosts(w http.ResponseWriter, _ *http.Request) {
	query := `SELECT * FROM posts ORDER BY id DESC`
	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		http.Error(w, "Failed to get posts from DB", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Body, &post.UserId); err != nil {
			http.Error(w, "Failed to scan post", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	CustomJsonResponse(w, http.StatusOK, posts)
}

func createPost(w http.ResponseWriter, r *http.Request) {
	var post Post

	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(post); err != nil {
		errs := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			switch err.Tag() {
			case "required":
				errs[field] = fmt.Sprintf("%s is required", field)
			case "min":
				errs[field] = fmt.Sprintf("%s must be at least %s characters", field, err.Param())
			}
		}

		CustomJsonResponse(w, http.StatusBadRequest, errs)
		return
	}

	query := `INSERT INTO posts (title, body, user_id) VALUES ($1, $2, $3) RETURNING id`

	var id int

	if err := DB.QueryRow(context.Background(), query, post.Title, post.Body, post.UserId).Scan(&id); err != nil {
		http.Error(w, "Failed to insert post in DB", http.StatusInternalServerError)
		return
	}

	post.ID = id

	CustomJsonResponse(w, http.StatusCreated, post)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("id")
	if postIdStr == "" {
		http.Error(w, "postIs is missing", http.StatusBadRequest)

		return
	}

	postId, err := strconv.Atoi(postIdStr)
	if err != nil || postId <= 0 {
		http.Error(w, "Invalid postId", http.StatusBadRequest)

		return
	}

	query := `DELETE FROM posts WHERE id = $1`

	result, err := DB.Exec(context.Background(), query, postId)
	if err != nil {
		http.Error(w, "Failed to delete post from DB", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	CustomJsonResponse(w, http.StatusOK, Success{Completed: true, Message: "Post deleted successfully"})
}

func updatePost(w http.ResponseWriter, r *http.Request) {
	var post Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Unable to get Request", http.StatusBadRequest)

		return
	}

	if err := validate.Struct(post); err != nil {
		errs := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()

			switch err.Tag() {
			case "required":
				errs[field] = fmt.Sprintf("%s is required", field)
			case "min":
				errs[field] = fmt.Sprintf("%s must be at least %s characters", field, err.Param())
			}
		}

		CustomJsonResponse(w, http.StatusBadRequest, errs)
		return
	}

	query := `UPDATE posts set title=$1, user_id=$2, body=$3 WHERE id=$4`
	result, err := DB.Exec(context.Background(), query, post.Title, post.UserId, post.Body, post.ID)
	if err != nil {
		http.Error(w, "Failed to update post from DB", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	CustomJsonResponse(w, http.StatusOK, post)

}
