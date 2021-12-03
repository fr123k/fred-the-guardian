package model

import (
	"fmt"
	"runtime"
)

const (
	INVALID_REQUEST_BODY = "E400"
	UNAUTHORIZED_REQUEST = "E401"
	TOO_MANY_REQUESTS    = "E429"
)

var (
	UNAUTHORIZED_REQUEST_RESPONSE = ErrorResponse{
		Code:    UNAUTHORIZED_REQUEST,
		Message: "Missing http header 'X-SECRET-KEY'.",
		Error:   "Unauthorized request",
	}

	INVALID_REQUEST_BODY_EMPTY_PAYLOAD = ErrorResponse{
		Code:    INVALID_REQUEST_BODY,
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

type Memory struct {
	Alloc uint64 `json:"alloc"`
	TotalAlloc uint64 `json:"total_alloc"`
	SysAlloc uint64 `json:"sys_alloc"`
	NumGC uint32 `json:"number_gc"`
}

type StatusResponse struct {
	Counters uint `json:"counters"`
	Memory Memory `json:"memory"`
}

func MemoryUsage() Memory {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return *InMegaBytes(&m)
}

func InMegaBytes (m *runtime.MemStats) *Memory {
	return &Memory{
		Alloc: BytesToMegaBytes(m.Alloc),
		TotalAlloc: BytesToMegaBytes(m.TotalAlloc),
		SysAlloc: BytesToMegaBytes(m.Sys),
		NumGC: m.NumGC,
	}
}

func BytesToMegaBytes(b uint64) uint64 {
    return b / 1024 / 1024
}

func TooManyRequests(maxCnt uint, wait int64) RateLimitResponse {
	return RateLimitResponse{
		ErrorResponse: ErrorResponse{
			Code:    TOO_MANY_REQUESTS,
			Message: fmt.Sprintf("Reach the limit of '%d' request.", maxCnt),
			Error:   "Rate limit exceeded."},
		Wait: wait,
	}
}

func InValidRequest(err string) ErrorResponse {
	return ErrorResponse{
		Code:    INVALID_REQUEST_BODY,
		Error:   err,
		Message: "Request body malformed.",
	}
}
