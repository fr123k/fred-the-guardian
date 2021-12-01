package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "net/http/httptest"
    "net/url"
    "os"
    "testing"

    "github.com/foxcpp/go-mockdns"
    "github.com/fr123k/fred-the-guardian/pkg/model"
    "github.com/stretchr/testify/assert"
)

func PingStub(t *testing.T) (*httptest.Server, *url.URL) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        switch r.RequestURI {
        case "/status":
            w.WriteHeader(http.StatusOK)
            return
        case "/ping":
            w.WriteHeader(http.StatusOK)
            w.Header().Set("Content-Type", "application/json")
            pong := model.PongResponse{
                Response: "pong",
            }
            json.NewEncoder(w).Encode(pong)
        default:
            http.Error(w, "not found", http.StatusNotFound)
            return
        }
    }))
    t.Cleanup(server.Close)

    url, err := url.Parse(server.URL)
    if err != nil {
        panic(err)
    }
    return server, url
}

func DNSStub(t *testing.T, u *url.URL) {
    // disable mockdns log output
    log := log.New(os.Stderr, "mockdns server: ", log.LstdFlags)
    log.SetOutput(ioutil.Discard)

    srv, _ := mockdns.NewServerWithLogger(map[string]mockdns.Zone{
        "fred.fr123k.uk.": {
            TXT: []string{"1 127.0.0.1:8888 /", fmt.Sprintf("5 %s /", u.Host)},
        },
    }, log, false)

    // defer srv.Close()
    t.Cleanup(func() {
        srv.Close()
    })

    srv.PatchNet(net.DefaultResolver)
    // defer mockdns.UnpatchNet(net.DefaultResolver)
    t.Cleanup(func() {
        mockdns.UnpatchNet(net.DefaultResolver)
    })
}

func TestAutodiscovery(t *testing.T) {

    _, u := PingStub(t)
    DNSStub(t, u)

    srvCfg := AutoDiscovery(PingConfig{secret: "Secret"})
    assert.Equal(t, "127.0.0.1", srvCfg.Server, "Server as expected")
    assert.Equal(t, u.Port(), srvCfg.Port, "Server as expected")
}

func TestAutodiscoveryFail(t *testing.T) {

    u := url.URL{
        Host: "127.0.0.1:36579",
    }
    DNSStub(t, &u)

    srvCfg := AutoDiscovery(PingConfig{secret: "Secret"})
    assert.Nil(t, srvCfg, "Autodiscovery failed so return nil.")
}

func DefaultPingConfig() PingConfig {
    return PingConfig{
        ServerConfig: ServerConfig{
            Server: DEFAULT_HOST,
            Port:   DEFAULT_PORT,
            path:   DEFAULT_PATH},
        secret:        DEFAULT_SECRET,
        AutoDiscovery: false,
        RandomSecret:  false,
    }
}

func TestParseArgs(t *testing.T) {

    pngConfig := parseArgs()

    assert.Equal(t, DEFAULT_HOST, pngConfig.Server, "Server as expected")
    assert.Equal(t, DEFAULT_PORT, pngConfig.Port, "Server as expected")
    assert.Equal(t, DEFAULT_PATH, pngConfig.Path(), "Server as expected")
    assert.Equal(t, DEFAULT_SECRET, pngConfig.Secret(), "Server as expected")
}

func TestPingClient(t *testing.T) {
    url, pngConfig := PingClient(DefaultPingConfig())
    assert.Equal(t, fmt.Sprintf("http://%s:%s%sping", DEFAULT_HOST, DEFAULT_PORT, DEFAULT_PATH), url, "Server as expected")
    assert.Equal(t, DEFAULT_HOST, pngConfig.Server, "Server as expected")
}

func TestPing(t *testing.T) {

    _, u := PingStub(t)
    DNSStub(t, u)

    pngCfg := DefaultPingConfig()
    pngCfg.AutoDiscovery = true
    url, pngConfig := PingClient(pngCfg)
    DoPingRequest(url, pngConfig)
    srvCfg := AutoDiscovery(PingConfig{secret: "Secret"})
    assert.Equal(t, "127.0.0.1", srvCfg.Server, "Server as expected")
    assert.Equal(t, u.Port(), srvCfg.Port, "Server as expected")
}
