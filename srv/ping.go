package main

import (
	"encoding/json"
	_"fmt"
	"net/http"
	"time"

	"github.com/fr123k/fred-the-guardian/pkg/counter"
	"github.com/fr123k/fred-the-guardian/pkg/middleware"
	"github.com/fr123k/fred-the-guardian/pkg/model"
	"github.com/fr123k/fred-the-guardian/pkg/utility"

	prommiddleware "github.com/albertogviana/prometheus-middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	bckCnt counter.Bucket
	rateLimitCntTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ratelimit_counter_total",
		Help: "Number of buckets with a counter.",
	})
)

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		rateLimitCntTotal.Set(float64(bckCnt.Size()))
	})
}

type HandlerFunc = func(w http.ResponseWriter, r *http.Request)

func status(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(model.StatusResponse{
		Counters: uint(bckCnt.Size()),
		Memory:   model.MemoryUsage(),
	})
}

func ping() HandlerFunc {
	validate := validator.New()
	validate.RegisterTagNameFunc(utility.JsonTagName)
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func main() {
	startPingService()
}

func startPingService() {
	router := startRouter()
	router = enablePrometheus(router)
	router = enableGlobalRateLimit(router)
	// router = enableBucketRateLimit(router)
	router = enableRedisBucketRateLimit(router)

	http.ListenAndServe(":"+utility.Env("PORT", "8080"), router)
}

func startRouter() *mux.Router {
	bckCnt = counter.NewBucketWitnCleanup(1 * time.Minute)
	// host     := utility.Env("REDIS_HOST", "localhost")
	// port     := utility.Env("REDIS_PORT", "6379")
	// password := utility.Env("REDIS_PASSWORD", "")

	// bckCnt = counter.NewRedisBucket(fmt.Sprintf("%s:%s", host, port), 1 * time.Minute)

	router := mux.NewRouter()

	router.HandleFunc("/ping", ping()).
		Methods("POST")
	router.HandleFunc("/status", status).
		Methods("GET")
	
	router.Use(middleware.AuthenticationMiddleware)

	return router
}

func enablePrometheus(router *mux.Router) *mux.Router {
	router.Path("/metrics").Handler(promhttp.Handler())
	router.Use(prommiddleware.NewPrometheusMiddleware(prommiddleware.Opts{}).InstrumentHandlerDuration)
	router.Use(prometheusMiddleware)
	return router
}

func enableGlobalRateLimit(router *mux.Router) *mux.Router {
	router.Use(middleware.GlobalCounterMiddleware(2, 1*time.Second))
	return router
}

func enableBucketRateLimit(router *mux.Router) *mux.Router {
	router.Use(middleware.BucketCountersMiddleware(&bckCnt, middleware.HTTP_HEADER_SECRET_KEY, 10))
	return router
}

func enableRedisBucketRateLimit(router *mux.Router) *mux.Router {
	router.Use(middleware.BucketCountersMiddleware(&bckCnt, middleware.HTTP_HEADER_SECRET_KEY, 10))
	return router
}
