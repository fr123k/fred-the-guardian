package main

import (
    "encoding/json"
    _ "encoding/json"
    "net/http"
    _ "strconv"

    "github.com/fr123k/fred-the-guardian/pkg/utility"
    "github.com/gorilla/mux"
)

type PongResponse struct {
    Response string `json:"response"`
}

func ping(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // params := mux.Vars(r)
    // total, err := strconv.ParseInt(params["total"], 10, 64)
    pong := PongResponse{
        Response: "pong",
    }
    json.NewEncoder(w).Encode(pong)
    return
}

// Middleware function, which will be called for each request
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _ = r.Header.Get("Authorization")
        // will trigger request processing
        next.ServeHTTP(w, r)
        // will stop request processing
        return
    })
}

func main() {

    router := mux.NewRouter()

    router.HandleFunc("/ping", ping).
        Methods("POST")

    router.Use(Middleware)

    http.ListenAndServe(":"+utility.Env("PORT", "8080"), router)
}
