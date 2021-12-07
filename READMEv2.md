<div id="top"></div>

<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![Build][build-shield]][build-url]
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/fr123k/fred-the-guardian">
    <img src="https://resume.fr123k.uk/images/logo_fr123k_transparent.png" alt="Logo" width="100" height="32">
  </a>

  <h3 align="center">Fred the Guardian</h3>

  <p align="center">
    <br />
    <a href="https://github.com/fr123k/fred-the-guardian"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/fr123k/fred-the-guardian">View Demo</a>
    <a href="https://github.com/fr123k/fred-the-guardian/issues">Report Bug</a>
    <a href="https://github.com/fr123k/fred-the-guardian/issues">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
        <li><a href="#command-line-interface-(pong)">Command Line Interface (pong)</a></li>
        <li><a href="#service-(ping)">Service (ping)</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li>
      <a href="#documentation">Documentation</a>
      <ul>
        <li><a href="#counter">Counter</a></li>
        <li><a href="#bucket">Bucket</a></li>
        <li><a href="#random-replacement">Random Replacement</a></li>
      </ul>
    </li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

<!--[![Product Name Screen Shot][product-screenshot]](https://example.com) -->
Implement rate limiting to study the the go programming language and learn things like random replacement, consistent hashing, leaky bucket, ... .


<p align="right">(<a href="#top">back to top</a>)</p>


### Built With

The project is written entirely in Golang.
* [go](https://go.dev/)

The gorilla mux framework is used for http routing and implementing the handlers that enforce rate limiting before a request is forwarded to the service.
* [gorilla/mux](https://github.com/gorilla/mux)

The package validator is used to validate the service's json request using struct's and field tags.
* [go-playground/validator](https://github.com/go-playground/validator/)

For monitoring, the prometheus client package is used to expose the service metrics.
* [prometheus](https://prometheus.io/)

The package prometheus-middleware defines the two important prometheus metrics for http based services. These are the latency and the number of requests per status code, method and path.
 * [albertogviana/prometheus-middleware](github.com/albertogviana/prometheus-middleware)

### Test With

The following packages are used only in the unit tests.

* [mcuadros/go-defaults](https://github.com/mcuadros/go-defaults)
* [foxcpp/go-mockdns](https://github.com/foxcpp/go-mockdns
* [stretchr/testify](https://github.com/stretchr/testify)


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- GETTING STARTED -->
## Getting Started

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

### Prerequisites

The following packages/tools has to be available.

* [make](https://www.gnu.org/software/make/)
* [go](https://go.dev/)
* [docker](https://www.docker.com/)
* [vagrant](https://www.vagrantup.com/)

### Installation

_Below is an example of how you can instruct your audience on installing and setting up your app. This template doesn't rely on any external dependencies or services._

1. Clone this repo
   ```bash
   git clone git@github.com:fr123k/fred-the-guardian.git
   #or
   git clone https://github.com/fr123k/fred-the-guardian.
   cd fred-the-guardian.
   ```
### Command Line Interface (pong)

#### Build

```bash
# build the pong command line interface
make build-cli
```

#### Run

```bash
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

# run with a specific ipaddress and port defintion for the ping service endpoint
./build/pong -server 192.168.2.45 -port 8888 -path /proxy/
```

### Service (ping)

#### Build

```bash
# build the go binary
make build
# or build the go binary within a docker image
make docker-build
```

#### Run

```bash
# run the go binary
make run
# or start the ping service docker container in detach mode
make docker--run
```

#### Test
```bash
make test
# or if the default port isn't available
PORT=8888 make test
# or just run the curl command yourself adjust the port as needed
curl -X POST -H 'X-SECRET-KEY:top secret' -v http://localhost:8080/ping
```

#### Deployment

The deployment has only been tested with a local minikube vm so far.
The deployment files are located under 'deploy/local'.

There is a [Makefile] in the local deployment folder.

```bash
cd deploy/local
# The vagrant make target will download the fr123k/ubuntu21-minikube vagrant box from vagrant cloud and 
# provision it based on the VagrantFile.
# The vagrant and the setup target steps can be skipped if you have aleady a local kubernetes cluster.
export KUBECONFIG_SAVE=${KUBECONFIG}
export KUBECONFIG=$(pwd)/kubectl.config
EXTERNAL_IP=172.28.128.16 make vagrant setup
# this just run kubectl apply for the ping k8s manifest files
make deploy
# this will use curl to send the ping request to the k8s pod
# if you use a different k8s cluster then the vagrant minikube one then
# adjust the EXTERNAL_IP
EXTERNAL_IP=172.28.128.16 make test
export KUBECONFIG=${KUBECONFIG_SAVE}
```


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- USAGE EXAMPLES -->
## Usage

Use this space to show useful examples of how a project can be used. Additional screenshots, code examples and demos work well in this space. You may also link to more resources.

_For more examples, please refer to the [Documentation](https://example.com)_


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- Documentation -->
## Documentation

### Counter

### Bucket

### Random Replacement


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- ROADMAP -->
## Roadmap

### Version 1

* [x] Implementation of the first simple rate limiting based on a counter in memory.
   * [x] single thread safe counter
   * [x] simple rate limiting with counter reset
   * [x] data structure for multiple counters

### Version 2

 * [x] Implementing the ping web server API without rate limiting
 * [x] Implementing the pong client interface
 * [x] Setting up a local minikube for the first deployment

### Version 3

 * [x] Implement rate limiting in the ping service with counters stored in memory.
 * [x] add ping service test
   * [x] Implementation of edge case and happy path tests
   * [x] Implementation of rate limiting tests
 * [x] refactoring the ping service, reducing the lines of code and abstracting the general behavior.
 * [x] add pong cli test
 * [x] implement pong cli handling of rate limiting responses.
 * [x] refactoring pong cli, reducing lines of code and abstracting general behavior.
 * [x] Add bucket count and memory usage to status endpoint response
 * [x] add prometheus metrics instrumentation for monitoring
 * [x] Implement eviction of bucket counters to purge memory of unused counters
 * [ ] Update documentation to the latest implemented features.

### Version 4

 * [ ] Implement rate limiting with counters stored in redis.

### Version 5

 * [ ] implement rate limiting based on an in memory [token bucket algorithm] (https://en.wikipedia.org/wiki/Token_bucket)

### Version 6

 * [ ] implement the rate limiting based on [token bucket algorithm](https://en.wikipedia.org/wiki/Token_bucket) stored in redis

### Version 7

 * [ ] implement rate limiting based on [fixed window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c) in memory
 * [ ] implement rate limiting based on [sliding window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c) in memory

### Version 8

 * [ ] implement rate limiting based on [fixed window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c) stored in redis
 * [ ] implement rate limiting based on [sliding window counters](https://medium.com/figma-design/an-alternative-approach-to-rate-limiting-f8a06cf7c94c) stored in redis

See the [open issues](https://github.com/fr123k/fred-the-guardian/issues) for a full list of proposed features (and known issues).


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- CONTACT -->
## Contact

Project Link: [https://github.com/fr123k/fred-the-guardian](https://github.com/fr123k/fred-the-guardian)


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

The following articles and resources provide useful insights into rate limiting and learning the Go programming language.

* [Random replacement (RR)](https://en.wikipedia.org/wiki/Cache_replacement_policies#Random_replacement_(RR))
* [Rate_limiting](https://en.wikipedia.org/wiki/Rate_limiting)
* [Github API rate limiter in Redis](https://github.blog/2021-04-05-how-we-scaled-github-api-sharded-replicated-rate-limiter-redis/)
* [effective go](https://go.dev/doc/effective_go)


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[build-shield]: https://img.shields.io/travis/fr123k/fred-the-guardian?style=for-the-badge
[build-url]: https://app.travis-ci.com/fr123k/fred-the-guardian
[contributors-shield]: https://img.shields.io/github/contributors/fr123k/fred-the-guardian?style=for-the-badge
[contributors-url]: https://github.com/fr123k/fred-the-guardian/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/fr123k/fred-the-guardian?style=for-the-badge
[forks-url]: https://github.com/fr123k/fred-the-guardian/network/members
[stars-shield]: https://img.shields.io/github/stars/fr123k/fred-the-guardian?style=for-the-badge
[stars-url]: https://github.com/fr123k/fred-the-guardian/stargazers
[issues-shield]: https://img.shields.io/github/issues/fr123k/fred-the-guardian?style=for-the-badge
[issues-url]: https://github.com/fr123k/fred-the-guardian/issues
[license-shield]: https://img.shields.io/github/license/fr123k/fred-the-guardian?style=for-the-badge
[license-url]: https://github.com/fr123k/fred-the-guardian/blob/main/LICENSE
[product-screenshot]: images/screenshot.png
