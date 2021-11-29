package pingtest

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"

    "github.com/fr123k/fred-the-guardian/pkg/model"
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
)

type PingTestExpect struct {
    Code int
    Body *string
}

type PingTest struct {
    Body   *string
    Secret *string
    Method string
    Path   string
    Expect PingTestExpect
}

func StringPtr(s string) *string {
    return &s
}

func (tp PingTest) Body_() io.Reader {
    if tp.Body == nil {
        return nil
    }
    return strings.NewReader(*tp.Body)
}

func (tp PingTest) Secret_() string {
    if tp.Secret != nil {
        return *tp.Secret
    }
    panic("Expect not nil string.")
}

func ExpectMissingHttpHeader() string {
    str, _ := json.Marshal(model.UNAUTHORIZED_REQUEST_RESPONSE)
    return fmt.Sprintf("%s\n", string(str))
}

func ExpectMissingProperPayload() string {
    str, _ := json.Marshal(model.INVALID_REQUEST_BODY_EMPTY_PAYLOAD)
    return fmt.Sprintf("%s\n", string(str))
}

func PingPayloadRequest(request string) string {
    str, _ := json.Marshal(model.PingRequest{
        Request: request,
    })
    return fmt.Sprintf("%s\n", string(str))
}

func PingPayloadResponse(response string) string {
    str, _ := json.Marshal(model.PongResponse{
        Response: response,
    })
    return fmt.Sprintf("%s\n", string(str))
}

func RunTests(tests map[string]PingTest, rt *mux.Router, t *testing.T) {
    for name, tc := range tests {
        t.Run(name, func(t *testing.T) {
            r, _ := http.NewRequest(tc.Method, tc.Path, tc.Body_())
            if tc.Secret != nil {
                r.Header.Set("X-SECRET-KEY", tc.Secret_())
            }
            w := httptest.NewRecorder()

            rt.ServeHTTP(w, r)

            assert.Equal(t, tc.Expect.Code, w.Code, name)
            assert.Equal(t, *tc.Expect.Body, w.Body.String())
        })
    }
}
