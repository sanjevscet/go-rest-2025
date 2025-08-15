package main

import (
	"fmt"
	"net/http"
	"time"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the Dummy API!")
}

func GetTimeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Current server time is:", time.Now())
}

func GetIPHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Client IP address is sanjeev:", r.RemoteAddr)
}
