package model


const (
    INVALID_REQUEST_BODY = "E400"
    UNAUTHORIZED_REQUEST = "E401"
    TOO_MANY_REQUESTS    = "E429"
)

var (
	UNAUTHORIZED_REQUEST_RESPONSE = ErrorResponse{
		Code: UNAUTHORIZED_REQUEST,
		Message: "Missing http header 'X-SECRET-KEY'.",
		Error: "Unauthorized request",
	}

	INVALID_REQUEST_BODY_EMPTY_PAYLOAD = ErrorResponse{
		Code: INVALID_REQUEST_BODY,
		Error:   "Missing proper payload",
		Message: "Request body malformed.",
	}
)

type PongResponse struct {
    Response string `json:"response"`
}
type PingRequest struct {
    Request string `json:"request" validate:"required"`
}
type ErrorResponse struct {
    Message string `json:"message"`
    Error   string `json:"error"`
    Code    string `json:"code"`
}

type RateLimitResponse struct {
    ErrorResponse
	Wait int64 `json:"wait"`
}
