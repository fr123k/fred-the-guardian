package main

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/fr123k/fred-the-guardian/pkg/pingtest"
)

func TestStatus(t *testing.T) {
	tests := []TestCase{
		PingTest{
			Name:   "ping request",
			Body:   StringPtr(PingPayloadRequest("ping")),
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusOK,
				Body: StringPtr(PingPayloadResponse("pong")),
			},
		},
		StatusTest{
			Name: "missing secret key",
			Body: nil,
			Expect: HttpTestExpect{
				Code: http.StatusUnauthorized,
				Body: StringPtr(ExpectMissingHttpHeader()),
			},
		},
		StatusTest{
			Name:   "status request with secret key",
			Body:   nil,
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusOK,
				Body: StringPtr(StatusResponse(1)),
			},
		},
	}
	router := startRouter()
	router = enableBucketRateLimit(router)
	RunHttpTests(tests, router, t)
}

func TestPing(t *testing.T) {
	tests := []TestCase{
		PingTest{
			Name: "missing secret key",
			Body: nil,
			Expect: HttpTestExpect{
				Code: http.StatusUnauthorized,
				Body: StringPtr(ExpectMissingHttpHeader()),
			},
		},
		PingTest{
			Name:   "empty payload",
			Body:   nil,
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusBadRequest,
				Body: StringPtr(ExpectMissingProperPayload()),
			},
		},
		PingTest{
			Name:   "non json payload",
			Body:   StringPtr("thats not a json payload"),
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusBadRequest,
				Body: StringPtr(ExpectInvalidRequest("invalid character 'h' in literal true (expecting 'r')")),
			},
		},
		PingTest{
			Name:   "malformed payload",
			Body:   StringPtr(MalformedBodyRequest("malformed body")),
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusBadRequest,
				Body: StringPtr(ExpectInvalidRequest("Key: 'PingRequest.request' Error:Field validation for 'request' failed on the 'required' tag")),
			},
		},
		PingTest{
			Name:   "ping request",
			Body:   StringPtr(PingPayloadRequest("ping")),
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusOK,
				Body: StringPtr(PingPayloadResponse("pong")),
			},
		},
		PingTest{
			Name:   "foobar request",
			Body:   StringPtr(PingPayloadRequest("foobar")),
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusOK,
				Body: StringPtr(PingPayloadResponse("foobar")),
			},
		},
	}

	RunHttpTests(tests, startRouter(), t)
}

func TestPingWithGlobalRateLimit(t *testing.T) {
	tests := []TestCase{
		PingTest{
			Name:   "ping request",
			Body:   StringPtr(PingPayloadRequest("ping")),
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusOK,
				Body: StringPtr(PingPayloadResponse("pong")),
			},
		},
		PingTest{
			Name:   "foobar request",
			Body:   StringPtr(PingPayloadRequest("foobar")),
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusOK,
				Body: StringPtr(PingPayloadResponse("foobar")),
			},
		},
		PingTest{
			Name:   "rate limit exceeded",
			Body:   StringPtr(PingPayloadRequest("foobar")),
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusTooManyRequests,
				Body: StringPtr(ExpectRateLimit(2, 1)),
			},
		},
	}
	router := startRouter()
	RunHttpTests(tests, enableGlobalRateLimit(router), t)
}

func multipleRequests(cnt uint) []TestCase {
	tests := make([]TestCase, cnt)
	for i := 0; i < 10; i++ {
		tests[i] = PingTest{
			Name:   fmt.Sprintf("ping request %d", i),
			Body:   StringPtr(PingPayloadRequest("ping")),
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusOK,
				Body: StringPtr(PingPayloadResponse("pong")),
			},
		}
	}
	return tests
}

func TestPingWithBucketRateLimit(t *testing.T) {
	tests := multipleRequests(10)
	rateLimitExceeded := []TestCase{
		PingTest{
			Name:   "rate limit exceeded",
			Body:   StringPtr(PingPayloadRequest("foobar")),
			Secret: StringPtr("secret"),
			Expect: HttpTestExpect{
				Code: http.StatusTooManyRequests,
				Body: StringPtr(ExpectRateLimit(10, 60)),
			},
		},
	}
	router := startRouter()
	router = enableBucketRateLimit(router)
	RunHttpTests(tests, router, t)
	RunHttpTests(rateLimitExceeded, router, t)
}
