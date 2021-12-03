package dns

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/foxcpp/go-mockdns"
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

func PingStub(t *testing.T) (*httptest.Server, *url.URL) {
	return HttpServerStub(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.RequestURI {
		default:
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
	})
}

func DNSStub(t *testing.T, port uint16) {
	// disable mockdns log output
	log := log.New(os.Stderr, "mockdns server: ", log.LstdFlags)
	log.SetOutput(ioutil.Discard)

	srv, _ := mockdns.NewServerWithLogger(map[string]mockdns.Zone{
		"localhost.localdomain.": {
			A: []string{"127.0.0.1"},
		},
		"fred.fr123k.uk.": {
			A: []string{"172.34.2.6"},
		},
		"_fred._tcp.fr123k.uk.": {
			SRV: []net.SRV{{
				Target:   "localhost.localdomain.",
				Port:     port,
				Priority: 10,
				Weight:   10,
			},
				{
					Target:   "fred.fr123k.uk.",
					Port:     port,
					Priority: 1,
					Weight:   10,
				},
			},
		},
		"_barney._tcp.fr123k.uk.": {
			SRV: []net.SRV{{
				Target:   "unknown.dns.domain.",
				Port:     port,
				Priority: 10,
				Weight:   10,
			},
			},
		},
	}, log, false)

	t.Cleanup(func() {
		mockdns.UnpatchNet(net.DefaultResolver)
		srv.Close()
	})

	srv.PatchNet(net.DefaultResolver)
}

func TestServiceDiscovery(t *testing.T) {
	_, u := PingStub(t)
	p, err := strconv.Atoi(u.Port())
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	DNSStub(t, uint16(p))

	srv := ServiceDiscovery("fred", TCPCheck)
	assert.Equal(t, Service{IP: "127.0.0.1", Host: "localhost.localdomain.", Port: uint16(p)}, *srv, "The localhost service has to be discovered from the possible options.")
}

func TestServiceDiscoveryFailNotReachableService(t *testing.T) {
	DNSStub(t, 34567)

	srv := ServiceDiscovery("fred", TCPCheck)
	assert.Nil(t, srv, "Service discovery returns nil if it doesn't found a valid DNS SRV record.")
}

func TestServiceDiscoveryFailUnknownServiceDomain(t *testing.T) {
	DNSStub(t, 34567)

	srv := ServiceDiscovery("barney", TCPCheck)
	assert.Nil(t, srv, "Service discovery returns nil if it doesn't found a valid DNS SRV record.")
}
