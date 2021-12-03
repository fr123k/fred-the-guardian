package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fr123k/fred-the-guardian/pkg/counter"
	"github.com/fr123k/fred-the-guardian/pkg/model"
	"github.com/gorilla/mux"
)

const (
	HTTP_HEADER_SECRET_KEY = "X-SECRET-KEY"
)

// Middleware function, which will be called for each request
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secKey := r.Header.Get(HTTP_HEADER_SECRET_KEY)

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
func BucketCountersMiddleware(counter *counter.Bucket, header string, maxCnt uint) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// will trigger request processing
			value := r.Header.Get(header)
			rate := counter.Increment(value)
			if rate.Count > uint64(maxCnt) {
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(model.TooManyRequests(maxCnt, rate.NextReset))
				// will stop request processing
				return
			}
			log.Printf("Bucket Rate %v, %d", rate, counter.Size())
			h.ServeHTTP(w, r)
			return
		})
	}
}
