package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/fr123k/fred-the-guardian/pkg/model"
)

const (
    DEFAULT_HOST   = "127.0.0.1"
    DEFAULT_PORT   = "8080"
    DEFAULT_PATH   = "/"
    DEFAULT_SECRET = "top secret"
)

type ServerConfig struct {
    Server string
    Port   string
    path   string
}

type PingConfig struct {
    AutoDiscovery bool
    RandomSecret  bool
    secret        string
    ServerConfig
}

func main() {
    ping()
}

func (p PingConfig) Host() string {
    return fmt.Sprintf("%s:%s", p.Server, p.Port)
}

func TrailingSlash(s string) string {
    if strings.HasSuffix(s, "/") {
        return s
    }
    return fmt.Sprintf("%s/", s)
}

func (p PingConfig) Path() string {
    return TrailingSlash(p.path)
}

//https://stackoverflow.com/a/50581165
func rndStr(len int) string {
    buff := make([]byte, len)
    rand.Read(buff)
    str := base64.StdEncoding.EncodeToString(buff)
    // Base 64 can be longer than len
    return str[:len]
}

func (p PingConfig) Secret() string {
    if p.RandomSecret {
        return rndStr(12)
    }
    return p.secret
}

func prettyPrint(i interface{}) string {
    s, _ := json.MarshalIndent(i, "", "\t")
    return string(s)
}

func SplitIntoTwoVars(str string, sep string) (string, string, error) {
    s := strings.Split(str, sep)
    if len(s) < 3 {
        return "", "", errors.New("Minimum match not found")
    }
    return s[1], s[2], nil
}

func autoDiscovery(pingCfg PingConfig) *ServerConfig {
    txtrecords, _ := net.LookupTXT("fred.fr123k.uk")

    sort.Strings(txtrecords)
    for _, txt := range txtrecords {
        host, path, _ := SplitIntoTwoVars(txt, " ")
        fmt.Printf("Auto Discovery try %s%s\n", host, path)
        client := &http.Client{}
        req, err := http.NewRequest("GET", fmt.Sprintf("http://%s%sstatus", host, path), nil)
        if err != nil {
            panic(err)
        }
        req.Header.Set("X-SECRET-KEY", pingCfg.Secret())
        resp, err := client.Do(req)
        if err != nil {
            log.Printf(err.Error())
            continue
        }
        log.Println(resp.StatusCode)
        if resp.StatusCode == 200 {
            server, port, err := net.SplitHostPort(host)
            if err != nil {
                log.Fatalf(err.Error())
                continue
            }
            return &ServerConfig{
                Server: server,
                Port:   port,
                path:   path,
            }
        }
    }
    return nil
}

func parseArgs() PingConfig {
    server := flag.String("server", DEFAULT_HOST, fmt.Sprintf("server address of the ping service (Default: %s)", DEFAULT_HOST))
    port := flag.String("port", DEFAULT_PORT, fmt.Sprintf("port of the ping service (Default: %s)", DEFAULT_PORT))
    path := flag.String("path", DEFAULT_PATH, fmt.Sprintf("root path of the ping service (Default: %s)", DEFAULT_PATH))

    secret := flag.String("secret", DEFAULT_SECRET, fmt.Sprintf("specify the secret value for the X-SECRET-KEY http header (Default: %s)", DEFAULT_SECRET))
    rndSecret := flag.Bool("rndsec", false, fmt.Sprintf("set true to generate a random secret for each request (Default: %t)", false))

    autoDiscovery := flag.Bool("auto", false, fmt.Sprintf("use auto discovery of possible ping services (Default: %t)", false))

    flag.Parse()

    return PingConfig{
        ServerConfig: ServerConfig{
            Server: *server,
            Port:   *port,
            path:   *path},
        AutoDiscovery: *autoDiscovery,
        secret:        *secret,
        RandomSecret: *rndSecret,
    }
}

func ping() {

    pingCfg := parseArgs()

    if pingCfg.AutoDiscovery {
        serverCfg := autoDiscovery(pingCfg)
        if serverCfg != nil {
            pingCfg.ServerConfig = *serverCfg
        }
    }
    pingRequest := model.PingRequest{
        Request: "ping",
    }

    fmt.Printf("Body %s\n", prettyPrint(pingRequest))

    url := fmt.Sprintf("http://%s%sping", pingCfg.Host(), pingCfg.Path())
    for {
        fmt.Printf("Call fred %s\n", url)
        client := &http.Client{}
        payloadBuf := new(bytes.Buffer)

        json.NewEncoder(payloadBuf).Encode(model.PingRequest{
            Request: "ping",
        })

        req, err := http.NewRequest("POST", url, payloadBuf)
        if err != nil {
            panic(err)
        }
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-SECRET-KEY", pingCfg.Secret())

        resp, err := client.Do(req)
        if err != nil {
            log.Printf("Error calling Jenkins '%s'\n", err)
            return
        }
        bodyText, err := ioutil.ReadAll(resp.Body)
        log.Printf("response %s", string(bodyText))
        time.Sleep(1 * time.Second)
    }
}
