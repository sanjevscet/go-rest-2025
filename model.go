package main

type Post struct {
	Id     int    `json:"id"`
	Title  string `json:"title" validate:"required,min=3"`
	Body   string `json:"body"`
	UserId int    `json:"userId"`
}

type Success struct {
	Completed bool   `json:"completed"`
	Message   string `json:"message"`
}
