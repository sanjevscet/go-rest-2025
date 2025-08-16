package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Post struct {
	ID     int    `json:"id"`
	Title  string `json:"title" validate:"required,min=3"`
	Body   string `json:"body"`
	UserId int    `json:"userId"`
}

type Success struct {
	Completed bool   `json:"completed"`
	Message   string `json:"message"`
}

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" validate:"required,min=3"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,min=6"`
	IsActive  bool      `json:"isActive" validate:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`

	jwt.RegisteredClaims
}

type UserLogin struct {
	Username string `json:"username" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=6"`
}
