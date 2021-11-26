# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.12
ARG GO_VERSION=1.17

# Second stage: build the executable
FROM golang:${GO_VERSION}-alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Import the Certificate-Authority certificates for the app to be able to send
# requests to HTTPS endpoints.
# RUN apk add --no-cache ca-certificates

# Accept the version of the app that will be injected into the compiled
# executable.
ARG APP_VERSION=undefined

# Set the environment variables for the build command.
ENV CGO_ENABLED=0 GOFLAGS=-mod=vendor

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /ping

COPY go.mod ./
COPY go.sum ./
COPY vendor ./vendor
RUN ls -lha
#RUN go mod vendor

# Import the code from the first stage.
# COPY --from=code /infra-hook ./
COPY srv ./srv
COPY pkg ./pkg

# Build the executable to `/app`. Mark the build as statically linked and
# inject the version as a global variable.
RUN go build \
    -installsuffix 'static' \
    -ldflags "-X main.Version=${APP_VERSION}" \
    -o /app \
    ./srv/ping.go
RUN go test -v -timeout 60s ./...

# Final stage: the running container
FROM scratch AS final

# This variable should be pass in the ci build pipeline
ARG APP_VERSION=undefined
ARG GIT_COMMIT=undefined
ARG BUILD_DATE=undefined

LABEL org.label-schema.build-date="$BUILD_DATE"
LABEL org.label-schema.name="ping"
LABEL org.label-schema.description="Ping web service"
LABEL org.label-schema.vcs-url="https://github.com/fr123k/fred-the.guardian"
LABEL org.label-schema.vcs-ref="$GIT_COMMIT"
LABEL org.label-schema.version="$APP_VERSION"
LABEL org.label-schema.schema-version="1.0"
LABEL go-version="${GO_VERSION}"

# Declare the port on which the application will be run.
EXPOSE 8080

# Import the user and group files.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable frmo the second stage.
COPY --from=builder /app /app

# Run the container as an unprivileged user.
USER nobody:nobody

ENTRYPOINT ["/app"]
