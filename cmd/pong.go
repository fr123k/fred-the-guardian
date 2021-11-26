package main

import (
    "bytes"
    "encoding/json"
    "errors"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "net/http"
    "sort"
    "strings"
    "time"

    "github.com/fr123k/fred-the-guardian/pkg/model"
    _ "github.com/fr123k/golang-template/pkg/utility"
)

const (
    DEFAULT_HOST = "172.28.128.16"
)

func main() {
    ping()
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

func autoDiscovery() (string, string) {
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
        req.Header.Set("X-SECRET-KEY", "top secret")
        resp, err := client.Do(req)
        if err != nil {
            log.Printf(err.Error())
            continue
        }
        log.Println(resp.StatusCode)
        if resp.StatusCode == 200 {
            return host, path
        }
    }
    return DEFAULT_HOST, "/"
}

func ping() {

    host, path := autoDiscovery()
    pingRequest := model.PingRequest{
        Request: "ping",
    }

    fmt.Printf("Body %s\n", prettyPrint(pingRequest))

    url := fmt.Sprintf("http://%s%sping", host, path)
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
        req.Header.Set("X-SECRET-KEY", "top secret")
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
