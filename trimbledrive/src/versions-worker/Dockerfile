ARG GO_VERSION=1.20
ARG APP_VERSION=v0.0.0

############################
# STEP 1 build executable binary in a separate builder image
############################
FROM --platform=linux/amd64 golang:${GO_VERSION}-alpine AS builder

ARG APP_VERSION

# Run the command below to install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

# Set the current working directory inside the container
RUN mkdir -p /build
WORKDIR /build

# Copy required folders for the build process.
COPY ./src/ ./src/

# Download api module dependencies
WORKDIR /build/src/versions-worker
RUN go mod tidy
RUN go mod download

############################
# STEP 2.1. build in debug mode
############################
FROM builder AS builder-debug
ARG APP_VERSION
WORKDIR /build/src/versions-worker
RUN go build -tags=debug -cover -race -ldflags="-X 'trimble.com/tdrive/versions-worker/Version=$APP_VERSION-debug'" -coverpkg=$(go list -f '{{if not .Standard}}{{.ImportPath}}{{end}}' -deps . | grep trimble | sed "s/trimble.com\/common/..\/common-go/" | sed "s/trimble.com\/tdrive\/versions-worker/./" | paste -s -d , /dev/stdin) -o ../../.bin/ .


############################
# STEP 2.2. build in release mode
############################
FROM builder AS builder-release
ARG APP_VERSION
WORKDIR /build/src/versions-worker
RUN go build -ldflags="-X 'trimble.com/tdrive/versions-worker/Version=$APP_VERSION'" -o ../../.bin/ .


############################
# STEP 3.1. build a small executable debug image
############################

FROM --platform=linux/amd64 alpine:latest AS debug

ARG APP_VERSION

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

# Copy our executable and config file(s)
COPY --from=builder-debug /build/.bin/* /.bin/
COPY --from=builder-debug /build/src/versions-worker/config*.yaml /.bin/

ENV GO_ENVIRONMENT development
ENV GOCOVERDIR /var/log/coverage/$APP_VERSION

# Run the binary
WORKDIR /.bin
USER nonroot
CMD echo Creating GOCOVERDIR $GOCOVERDIR... && mkdir -p $GOCOVERDIR && rm -rf $GOCOVERDIR/* && (ls $GOCOVERDIR || true) && while true; do echo Starting...; ./versions-worker; done


############################
# STEP 3.2. build a small executable release image
############################

FROM --platform=linux/amd64 alpine:latest AS release

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

# Copy our executable and config file(s)
COPY --from=builder-release /build/.bin/* /.bin/
COPY --from=builder-release /build/src/versions-worker/config*.yaml /.bin/

ENV GO_ENVIRONMENT production

# Run the Go binary,
WORKDIR /.bin
USER nonroot
ENTRYPOINT ["./versions-worker"]
