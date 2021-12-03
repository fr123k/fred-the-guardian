[![Build Status](https://travis-ci.com/fr123k/fred-the-guardian.svg?branch=main)](https://app.travis-ci.com/fr123k/fred-the-guardian)

# Fred the Guardian

## Introduction

Writing a little ping pong service that implements rate limiting with the programming language **golang**.

## Requirements

### Web application
 * Runs in a docker container
 * implements the following throttling policy:
   * maximum of 10 “ping” requests are allowed per x-secret-key and per minute
   * maximum of 2 requests per second are allowed, regardless of the x-secret-key

### CLI tool
 * Sends 1 “ping” request per second
 * Stops sending “ping” requests when limit of allowed requests is reached
 * Starts sending “ping” requests again when throttling is expired


## API

### Endpoint /ping

#### Request

Payload example:

`{ "request": "ping" }`

Required headers:

`x-secret-key: str`

“x-secret-key”: a random string

#### Response

response payload (when not throttled):

`{ "response": "pong" }`

Expected response payload(when throttled):

`{ "message": "request throttled request",  "throttle_age": int }`

Response attributes:

“message” - friendly message explaining what is happening

“throttle_age” - the elapsed time(in seconds) since the throttling has been applied

## Administration

There a couple of web servers and ingress controllers out there they implement different rate limiting algorithm already.

* apache
* nginx (experience with)
* kong
* istio
* ...

In my opinion its better to use one of those above if rate limiting itself is not the core of the companies business.

## Deployment

The deployment is only tested with a local minikube VM.
The deployment files are located at 'deploy/local'.

There is a [Makefile] in the local deployment folder

```bash
cd deploy/local
# download the fr123k/ubuntu21-minikube vagrant box from vagrant cloud
# provision it based on the VagrantFile
# the two step can be skipped if you have a kubernetes cluster already
export KUBECONFIG_SAVE=${KUBECONFIG}
export KUBECONFIG=$(pwd)/kubectl.config
make vagrant setup
# this just run kubectl apply for the ping k8s manifest files
make deploy
# this will use curl to send the ping request to the k8s pod
# if use use a different k8s cluster then the vagrant minikube one then
# adjust the EXTERNAL_IP
EXTERNAL_IP=172.28.128.16 make test
export KUBECONFIG=${KUBECONFIG_SAVE}
```

## Run Cli

```bash
# build the pong command line interface
make build-cli
# show the pong cli argument options
./build/pong -help
Usage of ./build/pong:
  -auto
        use auto discovery of possible ping services (Default: false)
  -path string
        root path of the ping service (Default: /) (default "/")
  -port string
        port of the ping service (Default: 8080) (default "8080")
  -rndsec
        set true to generate a random secret for each request (Default: false)
  -secret string
        specify the secret value for the X-SECRET-KEY http header (Default: top secret) (default "top secret")
  -server string
        server address of the ping service (Default: 127.0.0.1) (default "127.0.0.1")


# run with default ping service endpoint localhost:8080/
./build/pong

# run with auto discovery try couple of options to reach a ping service
./build/pong -auto true

# run with default ping service endpoint localhost:8080/
./build/pong -server 192.168.2.45 -port 8888 -path /proxy/

```

## Development

This would be the road map for a self developed rate limiting service.

### To do

#### Version 1

 * [x] implement the first simple rate limiting based on in memory counter
   * [x]  single thread safe counter
   * [x]  simple rate limiting with counter reset
   * [x]  data structure for multiple counters

#### Version 2

 * [x] implement the ping web server API without rate limiting
 * [x] implement the pong client interface
 * [x] setup local minikube for first deployment

#### Version 3

 * [x] implement rate limiting in the ping service with in memory stored counters
 * [x] add ping service test
   * [x] implement the edge case and happy path tests
   * [x] implement the rate policy tests
 * [x] refactor the ping service reduce code lines and abstract general behavior
 * [x] add pong cli test
 * [x] implement pong cli handling of rate limit responses
 * [x] refactor the pong cli reduce code lines and abstract general behavior
 * [x] add number of bucket counters and memory usage to the status endpoint response

#### Version 4

 * [ ] implement the rate limiting with counters stored in redis

#### Version 5

 * [ ] implement the rate limiting based on in memory [token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket)

#### Version 6

 * [ ] implement the rate limiting based on [token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket) stored in redis

#### Version 7

 * [ ] implement the rate limiting based on in memory [fixed window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c)
 * [ ] implement the rate limiting based on in memory [sliding window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c)

#### Version 8

 * [ ] implement the rate limiting based on [fixed window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c) stored in redis
 * [ ] implement the rate limiting based on [sliding window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c) stored in redis
