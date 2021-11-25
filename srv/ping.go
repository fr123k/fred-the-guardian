package main

import (
	_ "encoding/json"
	"fmt"
	"net/http"
	_ "strconv"

	"github.com/fr123k/fred-the-guardian/pkg/utility"
	"github.com/gorilla/mux"
)

func ping(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // params := mux.Vars(r)
    // total, err := strconv.ParseInt(params["total"], 10, 64)
    // json.NewEncoder(w).Encode(vmGroup)
    fmt.Println("pong")
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
        Methods("GET")

    router.Use(Middleware)

    http.ListenAndServe(":"+utility.Env("PORT", "8080"), router)
}
