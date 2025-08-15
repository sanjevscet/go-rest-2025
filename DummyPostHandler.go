package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DummyPostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getDummyPost(w, r)
	case http.MethodPost:
		createDummyPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getDummyPost(w http.ResponseWriter, r *http.Request) {
	postIdStr := r.URL.Query().Get("id")
	if postIdStr == "" {
		getAllPosts(w, r)
		return
	}

	postId, err := strconv.Atoi(postIdStr)
	if err != nil || postId <= 0 {
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://jsonplaceholder.typicode.com/posts/%d", postId)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch post", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	var post Post

	if err := json.NewDecoder(resp.Body).Decode(&post); err != nil {
		http.Error(w, "failed to decode response", http.StatusInternalServerError)
		return
	}

	CustomJsonResponse(w, http.StatusOK, post)
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	url := "https://jsonplaceholder.typicode.com/posts"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "unable to get all posts", http.StatusInternalServerError)
		return

	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(resp.Body)

	_, err = decoder.Token()

	if err != nil {
		http.Error(w, "failed to decode response", http.StatusInternalServerError)
		return
	}

	var posts []Post

	for decoder.More() {
		var post Post
		if err := decoder.Decode(&post); err != nil {
			http.Error(w, "failed to decode response", http.StatusInternalServerError)
			return
		}
		posts = append(posts, post)
	}

	CustomJsonResponse(w, http.StatusOK, posts)
}

func createDummyPost(w http.ResponseWriter, r *http.Request) {
	var newPost Post
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(newPost); err != nil {
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

	jsonData, err := json.Marshal(newPost)
	if err != nil {
		http.Error(w, "failed to encode the post", http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(http.MethodPost, "https://dummyjson.com/posts/add", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to call external API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var createdPost Post

	if err := json.NewDecoder(resp.Body).Decode(&createdPost); err != nil {
		http.Error(w, "failed to decode response", http.StatusInternalServerError)
		return
	}

	CustomJsonResponse(w, http.StatusCreated, createdPost)
}
