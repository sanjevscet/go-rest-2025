package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	InitDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/time", GetTimeHandler)
	mux.HandleFunc("/ip", GetIPHandler)
	mux.HandleFunc("/dummyPost", DummyPostHandler)
	mux.HandleFunc("/post", PostHandler)
	mux.HandleFunc("/user", UserHandler)

	LoggingMiddleware := LogMiddleware(mux)

	server := &http.Server{
		Addr:         ":1414",
		Handler:      LoggingMiddleware,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Custom server running on port %s\n", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
