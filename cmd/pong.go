package main

import (
    "bytes"
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
    "github.com/fr123k/fred-the-guardian/pkg/utility"
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

func (p PingConfig) Path() string {
    return utility.TrailingSlash(p.path)
}

func (p PingConfig) Secret() string {
    if p.RandomSecret {
        return utility.RandomString(12)
    }
    return p.secret
}

type DNSClient interface {
    LookupTXT(name string) ([]string, error)
}

func AutoDiscovery(pingCfg PingConfig) *ServerConfig {
    txtrecords, _ := net.LookupTXT("fred.fr123k.uk")

    sort.Strings(txtrecords)
    for _, txt := range txtrecords {
        host, path, _ := utility.SplitIntoTwoVars(txt, " ")
        log.Printf("Auto Discovery try %s%s\n", host, path)
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
        RandomSecret:  *rndSecret,
    }
}

type WaitFnc = func(time.Duration)

func DoPingRequest(url string, pingCfg PingConfig, wait WaitFnc) {
    log.Printf("Call fred %s\n", url)
    client := &http.Client{}
    payloadBuf := new(bytes.Buffer)

    pingRequest := model.PingRequest{
        Request: "ping",
    }

    log.Printf("Request: %s\n", utility.PrettyPrint(pingRequest))
    json.NewEncoder(payloadBuf).Encode(pingRequest)

    req, err := http.NewRequest("POST", url, payloadBuf)
    if err != nil {
        panic(err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-SECRET-KEY", pingCfg.Secret())

    r, err := client.Do(req)
    if err != nil {
        log.Printf("Error calling fred '%s'\n", err)
        return
    }
    
    switch r.StatusCode {
        case http.StatusOK:
            var pongRsp model.PongResponse
            err = json.NewDecoder(r.Body).Decode(&pongRsp)
            if err != nil {
                log.Printf("Error json decoding '%s'\n", err)
                return
            }
            defer r.Body.Close()
            log.Printf("Response: %v", utility.PrettyPrint(pongRsp))

        case http.StatusTooManyRequests:
            var rateRsp model.RateLimitResponse
            err = json.NewDecoder(r.Body).Decode(&rateRsp)
            if err != nil {
                log.Printf("Error json decoding '%s'\n", err)
                return
            }
            log.Printf("Response: %v", utility.PrettyPrint(rateRsp))
            wait(time.Duration(rateRsp.Wait) * time.Second)
            defer r.Body.Close()
        default:
            bodyText, _ := ioutil.ReadAll(r.Body)
            log.Printf("Response: %s", string(bodyText))
    }
}

func PingClient(pingCfg PingConfig) (string, PingConfig) {
    if pingCfg.AutoDiscovery {
        serverCfg := AutoDiscovery(pingCfg)
        if serverCfg != nil {
            pingCfg.ServerConfig = *serverCfg
        }
    }

    url := fmt.Sprintf("http://%s%sping", pingCfg.Host(), pingCfg.Path())
    return url, pingCfg
}

func ping() {
    pingCfg := parseArgs()

    url, pingCfg := PingClient(pingCfg)
    for {
        DoPingRequest(url, pingCfg, func(duration time.Duration){
            log.Printf("Wait for %s that the rate limit can cool down.\n", duration)
            time.Sleep(duration)
        })
        time.Sleep(1 * time.Second)
    }
}
