package main

import (
    "fmt"
    "net/http"
    "github.com/afonsocraposo/go-rate-limiter/internal/handlers"
    "github.com/afonsocraposo/go-rate-limiter/internal/rate-limiter"
)

func main() {
    http.HandleFunc("/limited", ratelimiter.RateLimiter(handlers.Limited, 10, 1000/1))
    http.HandleFunc("/unlimited", handlers.Unlimited)
    fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

