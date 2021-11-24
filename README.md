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

## Development

This would be the road map for a self developed rate limiting service.

### To do

#### Version 1

 * [ ] implement the first simple rate limiting based on in memory counter
   * [x]  single thread safe counter
   * [x]  simple rate limiting with counter reset
   * [ ]  data structure for multiple counters

#### Version 2

 * [ ] implement the ping web server API without rate limiting
 * [ ] implement the pong client interface
 * [ ] setup local minibike for first deployment

#### Version 3

 * [ ] implement the rate limiting with counters stored in redis

#### Version 4

 * [ ] implement the rate limiting based on in memory [token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket)

#### Version 5 

 * [ ] implement the rate limiting based on [token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket) stored in redis

#### Version 6

 * [ ] implement the rate limiting based on in memory [fixed window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c)
 * [ ] implement the rate limiting based on in memory [sliding window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c)

#### Version 7

 * [ ] implement the rate limiting based on [fixed window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c) stored in redis
 * [ ] implement the rate limiting based on [sliding window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c) stored in redis
