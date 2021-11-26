package main

import (
	"encoding/json"
	_ "encoding/json"
	"net/http"
	"reflect"
	_ "strconv"
	"strings"

	"github.com/fr123k/fred-the-guardian/pkg/utility"
	"github.com/fr123k/fred-the-guardian/pkg/model"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

const (
	INVALID_REQUEST_BODY = "E400"
	UNAUTHORIZED_REQUEST = "E401"
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
	// params := mux.Vars(r)
	// total, err := strconv.ParseInt(params["total"], 10, 64)

	err := json.NewDecoder(r.Body).Decode(&pingRqt)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{
			Code: INVALID_REQUEST_BODY,
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
		json.NewEncoder(w).Encode(model.ErrorResponse{
			Code:    INVALID_REQUEST_BODY,
			Error:   validationErrors.Error(),
			Message: "Request body malformed.",
		})
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
			json.NewEncoder(w).Encode(model.ErrorResponse{
				Code: UNAUTHORIZED_REQUEST,
				//TODO expose service internal error message is not good security practice but good for quick development
				Message: "Missing http header 'X-SECRET-KEY'.",
			})
			// will stop request processing
			return
		}
		// will trigger request processing
		next.ServeHTTP(w, r)
		return
	})
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/ping", ping).
		Methods("POST")
	router.HandleFunc("/status", status).
		Methods("GET")

	router.Use(Middleware)

	http.ListenAndServe(":"+utility.Env("PORT", "8080"), router)
}
