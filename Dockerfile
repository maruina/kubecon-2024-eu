# Build the tests
FROM golang:1.21 AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy the go source
COPY e2e/ e2e/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go test -c -o e2e.test ./e2e

# Run the tests
FROM gcr.io/distroless/static:nonroot
WORKDIR /

COPY --from=builder /workspace/e2e.test .
USER 65532:65532

ENTRYPOINT ["/e2e.test"]
