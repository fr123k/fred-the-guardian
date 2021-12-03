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
	"time"

	"github.com/foxcpp/go-mockdns"
	"github.com/fr123k/fred-the-guardian/pkg/model"
	"github.com/stretchr/testify/assert"
)

func HttpServerStub(t *testing.T, handlerFunc http.HandlerFunc) (*httptest.Server, *url.URL) {
	server := httptest.NewServer(handlerFunc)
	t.Cleanup(server.Close)

	url, err := url.Parse(server.URL)
	if err != nil {
		panic(err)
	}
	return server, url
}

func PingStubRateLimitExceeded(t *testing.T) (*httptest.Server, *url.URL) {
	return HttpServerStub(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/status":
			w.WriteHeader(http.StatusOK)
			return
		case "/ping":
			w.WriteHeader(http.StatusTooManyRequests)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(model.TooManyRequests(13, 16))
		default:
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
	})
}

func PingStub(t *testing.T) (*httptest.Server, *url.URL) {
	return HttpServerStub(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		case "/status":
			w.WriteHeader(http.StatusOK)
			return
		case "/ping":
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(model.PongResponse{
				Response: "pong",
			})
		default:
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
	})
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

	t.Cleanup(func() {
		mockdns.UnpatchNet(net.DefaultResolver)
		srv.Close()
	})

	srv.PatchNet(net.DefaultResolver)
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

func TestPong(t *testing.T) {
	_, u := PingStub(t)
	DNSStub(t, u)

	pngCfg := DefaultPingConfig()
	pngCfg.AutoDiscovery = true
	url, pngConfig := PingClient(pngCfg)
	DoPingRequest(url, pngConfig, func(duration time.Duration) {
		assert.FailNow(t, "The wait function is not expected to be called in this test.")
	})
	srvCfg := AutoDiscovery(PingConfig{secret: "Secret"})
	assert.Equal(t, "127.0.0.1", srvCfg.Server, "Server as expected")
	assert.Equal(t, u.Port(), srvCfg.Port, "Server as expected")
}

func TestPongRateLimitExceeded(t *testing.T) {
	_, u := PingStubRateLimitExceeded(t)
	DNSStub(t, u)

	pngCfg := DefaultPingConfig()
	pngCfg.AutoDiscovery = true
	url, pngConfig := PingClient(pngCfg)
	DoPingRequest(url, pngConfig, func(duration time.Duration) {
		assert.Equal(t, time.Duration(16)*time.Second, duration, "Server as expected")
	})
	srvCfg := AutoDiscovery(PingConfig{secret: "Secret"})
	assert.Equal(t, "127.0.0.1", srvCfg.Server, "Server as expected")
	assert.Equal(t, u.Port(), srvCfg.Port, "Server as expected")
}
