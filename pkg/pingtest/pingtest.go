package HttpTest

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
    "github.com/mcuadros/go-defaults"
    "github.com/stretchr/testify/assert"
)

type Test interface {
    toTest() HttpTest
}

type HttpTestExpect struct {
    Code int
    Body *string
}

type HttpTest struct {
    Body   *string
    Secret *string
    Method string
    Path   string
    Expect HttpTestExpect
}

type StatusTest struct {
    Body   *string
    Secret *string
    Method string `default:"GET"`
    Path   string `default:"/status"`
    Expect HttpTestExpect
}

type PingTest struct {
    Body   *string
    Secret *string
    Method string `default:"POST"`
    Path   string `default:"/ping"`
    Expect HttpTestExpect
}

func (p *PingTest) toTest() HttpTest {
    defaults.SetDefaults(p)
    return HttpTest{
        Body: p.Body,
        Secret: p.Secret,
        Method: p.Method,
        Path: p.Path,
        Expect: p.Expect,
    }
}

func (s *StatusTest) toTest() HttpTest {
    defaults.SetDefaults(s)
    return HttpTest{
        Body: s.Body,
        Secret: s.Secret,
        Method: s.Method,
        Path: s.Path,
        Expect: s.Expect,
    }
}

func StringPtr(s string) *string {
    return &s
}

func (tp HttpTest) Body_() io.Reader {
    if tp.Body == nil {
        return nil
    }
    return strings.NewReader(*tp.Body)
}

func (tp HttpTest) Secret_() string {
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

func ExpectInvalidRequest(err string) string {
    str, _ := json.Marshal(model.InValidRequest(err))
    return fmt.Sprintf("%s\n", string(str))
}

func ExpectRateLimit(maxCnt uint, wait int64) string {
    str, _ := json.Marshal(model.TooManyRequests(maxCnt, wait))
    return fmt.Sprintf("%s\n", string(str))
}

func PingPayloadRequest(request string) string {
    str, _ := json.Marshal(model.PingRequest{
        Request: request,
    })
    return fmt.Sprintf("%s\n", string(str))
}

type MalformedRequest struct {
    Unknown string `json:"unknown"`
}

func MalformedBodyRequest(request string) string {
    str, _ := json.Marshal(MalformedRequest{
        Unknown: request,
    })
    return fmt.Sprintf("%s\n", string(str))
}

func PingPayloadResponse(response string) string {
    str, _ := json.Marshal(model.PongResponse{
        Response: response,
    })
    return fmt.Sprintf("%s\n", string(str))
}

type toHttpTestFunc func(map[string]HttpTest) 

//TODO check if this boilerplate code can be reduced
func toHttpTest(size uint, f toHttpTestFunc) map[string]HttpTest {
    httpTests := make(map[string]HttpTest, size)
    f(httpTests)
    return httpTests
}

func PingToHttpTests(tests map[string]PingTest) map[string]HttpTest {
    return toHttpTest(uint(len(tests)), func(httpTests map[string]HttpTest) {
        for name, tc := range tests {
            httpTests[name] = tc.toTest()
        }
    })
}

func StatusToHttpTests(tests map[string]StatusTest) map[string]HttpTest {
    return toHttpTest(uint(len(tests)), func(httpTests map[string]HttpTest) {
        for name, tc := range tests {
            httpTests[name] = tc.toTest()
        }
    })
}

func RunHttpTests(tests map[string]HttpTest, rt *mux.Router, t *testing.T) {
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
