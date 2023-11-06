package main

import (
    "fmt"
    "net/http"
    "github.com/afonsocraposo/go-rate-limiter/internal"
)

func main() {
    http.HandleFunc("/limited", handlers.Limited)
    http.HandleFunc("/unlimited", handlers.Unlimited)
    fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

