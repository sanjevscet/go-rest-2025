package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		login(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	var userLogin UserLogin

	if err := json.NewDecoder(r.Body).Decode(&userLogin); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(userLogin); err != nil {
		errs := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			errs[err.Field()] = err.Tag()
		}
		http.Error(w, "Validation failed", http.StatusBadRequest)
		return
	}

	// Authenticate user
	if err := authenticateUser(userLogin); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateJWTWithClaims(User{
		Username: userLogin.Username,
	})
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	CustomJsonResponse(w, http.StatusOK, map[string]string{"token": token})

}

func authenticateUser(userLogin UserLogin) error {

	query := `SELECT id, username, password FROM users WHERE username = $1`
	var user User

	if err := DB.QueryRow(context.Background(), query, userLogin.Username).Scan(&user.ID, &user.Username, &user.Password); err != nil {
		return err
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		return err
	}

	return nil
}
