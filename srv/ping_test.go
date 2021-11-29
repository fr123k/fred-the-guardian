package main

import (
    "net/http"
    "testing"

    . "github.com/fr123k/fred-the-guardian/pkg/pingtest"
)

func TestStatus(t *testing.T) {
    tests := map[string]PingTest{
        "missing secret key": {
            Body:   nil,
            Path:   "/status",
            Method: "GET",
            Expect: PingTestExpect{
                Code: http.StatusUnauthorized,
                Body: StringPtr(ExpectMissingHttpHeader()),
            },
        },
        "secret key": {
            Body:   nil,
            Path:   "/status",
            Method: "GET",
            Secret: StringPtr("secret"),
            Expect: PingTestExpect{
                Code: http.StatusOK,
                Body: StringPtr(""),
            },
        },
    }

    RunTests(tests, startRouter(), t)
}

func TestPingEdgeCase(t *testing.T) {
    tests := map[string]PingTest{
        "missing secret key": {
            Body:   nil,
            Path:   "/ping",
            Method: "POST",
            Expect: PingTestExpect{
                Code: http.StatusUnauthorized,
                Body: StringPtr(ExpectMissingHttpHeader()),
            },
        },
        "empty payload": {
            Body:   nil,
            Path:   "/ping",
            Method: "POST",
            Secret: StringPtr("secret"),
            Expect: PingTestExpect{
                Code: http.StatusBadRequest,
                Body: StringPtr(ExpectMissingProperPayload()),
            },
        },
    }

    RunTests(tests, startRouter(), t)
}

func TestPing(t *testing.T) {
    tests := map[string]PingTest{
        "ping request": {
            Body:   StringPtr(PingPayloadRequest("ping")),
            Path:   "/ping",
            Method: "POST",
            Secret: StringPtr("secret"),
            Expect: PingTestExpect{
                Code: http.StatusOK,
                Body: StringPtr(PingPayloadResponse("pong")),
            },
        },
        "foobar request": {
            Body:   StringPtr(PingPayloadRequest("foobar")),
            Path:   "/ping",
            Method: "POST",
            Secret: StringPtr("secret"),
            Expect: PingTestExpect{
                Code: http.StatusOK,
                Body: StringPtr(PingPayloadResponse("foobar")),
            },
        },
    }

    RunTests(tests, startRouter(), t)
}
