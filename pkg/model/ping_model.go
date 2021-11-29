package model

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
