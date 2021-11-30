package main

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/fr123k/fred-the-guardian/pkg/pingtest"
)

func TestStatus(t *testing.T) {
    tests := map[string]StatusTest{
        "missing secret key": {
            Body:   nil,
            Expect: HttpTestExpect{
                Code: http.StatusUnauthorized,
                Body: StringPtr(ExpectMissingHttpHeader()),
            },
        },
        "secret key": {
            Body:   nil,
            Secret: StringPtr("secret"),
            Expect: HttpTestExpect{
                Code: http.StatusOK,
                Body: StringPtr(""),
            },
        },
    }

    RunHttpTests(StatusToHttpTests(tests), startRouter(), t)
}

func TestPing(t *testing.T) {
    tests := map[string]PingTest{
        "missing secret key": {
            Body: nil,
            Expect: HttpTestExpect{
                Code: http.StatusUnauthorized,
                Body: StringPtr(ExpectMissingHttpHeader()),
            },
        },
        "empty payload": {
            Body:   nil,
            Secret: StringPtr("secret"),
            Expect: HttpTestExpect{
                Code: http.StatusBadRequest,
                Body: StringPtr(ExpectMissingProperPayload()),
            },
        },
        "non json payload": {
            Body:   StringPtr("thats not a json payload"),
            Secret: StringPtr("secret"),
            Expect: HttpTestExpect{
                Code: http.StatusBadRequest,
                Body: StringPtr(ExpectInvalidRequest("invalid character 'h' in literal true (expecting 'r')")),
            },
        },
        "malformed payload": {
            Body:   StringPtr(MalformedBodyRequest("malformed body")),
            Secret: StringPtr("secret"),
            Expect: HttpTestExpect{
                Code: http.StatusBadRequest,
                Body: StringPtr(ExpectInvalidRequest("Key: 'PingRequest.request' Error:Field validation for 'request' failed on the 'required' tag")),
            },
        },
        "ping request": {
            Body:   StringPtr(PingPayloadRequest("ping")),
            Secret: StringPtr("secret"),
            Expect: HttpTestExpect{
                Code: http.StatusOK,
                Body: StringPtr(PingPayloadResponse("pong")),
            },
        },
        "foobar request": {
            Body:   StringPtr(PingPayloadRequest("foobar")),
            Secret: StringPtr("secret"),
            Expect: HttpTestExpect{
                Code: http.StatusOK,
                Body: StringPtr(PingPayloadResponse("foobar")),
            },
        },
    }

    RunHttpTests(PingToHttpTests(tests), startRouter(), t)
}

func TestPingWithGlobalRateLimit(t *testing.T) {
    tests := map[string]PingTest{
        "ping request": {
            Body:   StringPtr(PingPayloadRequest("ping")),
            Secret: StringPtr("secret"),
            Expect: HttpTestExpect{
                Code: http.StatusOK,
                Body: StringPtr(PingPayloadResponse("pong")),
            },
        },
        "foobar request": {
            Body:   StringPtr(PingPayloadRequest("foobar")),
            Secret: StringPtr("secret"),
            Expect: HttpTestExpect{
                Code: http.StatusOK,
                Body: StringPtr(PingPayloadResponse("foobar")),
            },
        },
        "rate limit exceeded": {
            Body:   StringPtr(PingPayloadRequest("foobar")),
            Secret: StringPtr("secret"),
            Expect: HttpTestExpect{
                Code: http.StatusTooManyRequests,
                Body: StringPtr(ExpectRateLimit(2, 1)),
            },
        },
    }
    router := startRouter()
    RunHttpTests(PingToHttpTests(tests), enableGlobalRateLimit(router), t)
}

func multipleRequests(cnt uint) map[string]PingTest {
    tests := map[string]PingTest{}
    for i := 0; i < 10; i++ {
        tests[fmt.Sprintf("ping request %d", i)] = PingTest{
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
    rateLimitExceeded := map[string]PingTest{
        "rate limit exceeded": {
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
    RunHttpTests(PingToHttpTests(tests), router, t)
    RunHttpTests(PingToHttpTests(rateLimitExceeded), router, t)

}
