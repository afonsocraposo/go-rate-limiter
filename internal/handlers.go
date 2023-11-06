package handlers

import (
    "fmt"
    "net/http"
)

func Limited(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Limited, don't over use me!")
}

func Unlimited(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, "Unlimited! Let's Go!")
}
