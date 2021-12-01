package main

import (
    "encoding/json"
    "log"
    "net/http"
    "reflect"
    "strings"
    "time"

    "github.com/fr123k/fred-the-guardian/pkg/counter"
    "github.com/fr123k/fred-the-guardian/pkg/model"
    "github.com/fr123k/fred-the-guardian/pkg/utility"

    "github.com/go-playground/validator/v10"
    "github.com/gorilla/mux"
)

// Incase of invalid struct it returns the field name from the json tag instead of the struct variable name.
func jsonTagName(fld reflect.StructField) string {
    name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

    if name == "-" {
        return ""
    }

    return name
}

func status(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

func ping(w http.ResponseWriter, r *http.Request) {
    var pingRqt model.PingRequest

    w.Header().Set("Content-Type", "application/json")

    if r.Body == nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(model.INVALID_REQUEST_BODY_EMPTY_PAYLOAD)
        return
    }

    err := json.NewDecoder(r.Body).Decode(&pingRqt)
    defer r.Body.Close()

    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(model.ErrorResponse{
            Code: model.INVALID_REQUEST_BODY,
            //TODO expose service internal error message is not good security practice but good for quick development
            Error:   err.Error(),
            Message: "Request body malformed.",
        })
        return
    }

    validate := validator.New()
    validate.RegisterTagNameFunc(jsonTagName)

    err = validate.Struct(pingRqt)
    if err != nil {
        validationErrors := err.(validator.ValidationErrors)
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(model.InValidRequest(validationErrors.Error()))
        return
    }

    var pong model.PongResponse

    if pingRqt.Request == "ping" {
        pong = model.PongResponse{
            Response: "pong",
        }
    } else {
        //TODO not defined in requirements. Returns the sended request string in the response message
        pong = model.PongResponse{
            Response: pingRqt.Request,
        }
    }

    json.NewEncoder(w).Encode(pong)
    return
}

// Middleware function, which will be called for each request
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        secKey := r.Header.Get("X-SECRET-KEY")

        if len(secKey) == 0 {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(model.UNAUTHORIZED_REQUEST_RESPONSE)
            // will stop request processing
            return
        }

        // will trigger request processing
        next.ServeHTTP(w, r)
        return
    })
}

// Middleware function, which will be called for each request
// TODO add name to identify it in logs and tests
func GlobalCounterMiddleware(maxCnt uint, duration time.Duration) mux.MiddlewareFunc {
    counter := counter.NewRateLimit(duration)
    return func(h http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // will trigger request processing
            rate := counter.Increment()
            if rate.Count > uint64(maxCnt) {
                w.WriteHeader(http.StatusTooManyRequests)
                json.NewEncoder(w).Encode(model.TooManyRequests(maxCnt, rate.NextReset))
                // will stop request processing
                return
            }
            log.Printf("Global Rate %v", rate)
            h.ServeHTTP(w, r)
            return
        })
    }
}

// Middleware function, which will be called for each request
// TODO add name to identify it in logs and tests
func BucketCountersMiddleware(key string, maxCnt uint, duration time.Duration) mux.MiddlewareFunc {
    counter := counter.NewBucket(duration)
    return func(h http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // will trigger request processing
            secKey := r.Header.Get(key)
            rate := counter.Increment(secKey)
            if rate.Count > uint64(maxCnt) {
                w.WriteHeader(http.StatusTooManyRequests)
                json.NewEncoder(w).Encode(model.TooManyRequests(maxCnt, rate.NextReset))
                // will stop request processing
                return
            }
            log.Printf("Bucket Rate %v", rate)
            h.ServeHTTP(w, r)
            return
        })
    }
}

func main() {
    router := startRouter()
    router = enableGlobalRateLimit(router)
    router = enableBucketRateLimit(router)
    http.ListenAndServe(":"+utility.Env("PORT", "8080"), router)
}

func startRouter() *mux.Router {
    router := mux.NewRouter()

    router.HandleFunc("/ping", ping).
        Methods("POST")
    router.HandleFunc("/status", status).
        Methods("GET")

    router.Use(Middleware)

    return router
}

func enableGlobalRateLimit(router *mux.Router) *mux.Router {
    router.Use(GlobalCounterMiddleware(2, 1 * time.Second))
    return router
}

func enableBucketRateLimit(router *mux.Router) *mux.Router {
    router.Use(BucketCountersMiddleware("X-SECRET-KEY", 10, 1 * time.Minute))
    return router
}
