package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUser(w, r)
	case http.MethodGet:
		getUser(w, r)
	case http.MethodPut:
		updateUser(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(user); err != nil {
		errs := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			switch err.Tag() {
			case "required":
				errs[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errs[field] = fmt.Sprintf("%s is not a valid email", field)
			case "min":
				errs[field] = fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
			}
		}

		CustomJsonResponse(w, http.StatusBadRequest, errs)
		return
	}

	if err := checkUsernameOrEmail(user); err != nil {
		CustomJsonResponse(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		CustomJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		return
	}

	query := `INSERT INTO users (username, email, password, is_active) VALUES ($1, $2, $3, $4) RETURNING id`
	var userID int
	if err := DB.QueryRow(context.Background(), query, user.Username, user.Email, hashedPassword, user.IsActive).Scan(&userID); err != nil {
		CustomJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		return
	}

	CustomJsonResponse(w, http.StatusCreated, map[string]int{"id": userID})
}

func checkUsernameOrEmail(user User) error {
	var count int

	args := []any{user.Username, user.Email}

	query := `SELECT COUNT(*) FROM users WHERE (username = $1 OR email = $2)`
	if user.ID != 0 {
		query += ` AND id != $3`
		args = append(args, user.ID)
	}
	if err := DB.QueryRow(context.Background(), query, args...).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("username or email already exists")
	}
	return nil
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		getUsers(w, r)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `SELECT id, username, email, is_active, created_at, updated_at FROM users WHERE id = $1`

	var user User

	if err := DB.QueryRow(context.Background(), query, userID).Scan(&user.ID, &user.Username, &user.Email, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {

		log.Println(err)

		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		return
	}

	user.ID = userID

	CustomJsonResponse(w, http.StatusOK, user)
}

func getUsers(w http.ResponseWriter, _ *http.Request) {
	users := make([]User, 0)

	query := `SELECT id, username, email, is_active, created_at, updated_at FROM users`
	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Println(err)
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	CustomJsonResponse(w, http.StatusOK, users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(user); err != nil {
		errs := make(map[string]string)

		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			switch err.Tag() {
			case "required":
				errs[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errs[field] = fmt.Sprintf("%s is not a valid email", field)
			case "min":
				errs[field] = fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
			}
		}

		CustomJsonResponse(w, http.StatusBadRequest, errs)
		return
	}

	if err := checkUsernameOrEmail(user); err != nil {
		CustomJsonResponse(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		CustomJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		return
	}

	query := `UPDATE users SET username = $1, email = $2, password = $3, is_active = $4, updated_at = $5 WHERE id = $6`
	if _, err := DB.Exec(context.Background(), query, user.Username, user.Email, hashedPassword, user.IsActive, time.Now(), user.ID); err != nil {
		CustomJsonResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
		return
	}

	CustomJsonResponse(w, http.StatusOK, map[string]string{"message": "User updated successfully"})
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM users WHERE id = $1`
	result, err := DB.Exec(context.Background(), query, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected() == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	CustomJsonResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
