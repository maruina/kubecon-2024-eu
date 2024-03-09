# Build the tests binary
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
COPY e2e-runner.sh .

# Build
RUN chmod +x e2e-runner.sh \
    && CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -o test2json -ldflags="-s -w" cmd/test2json \
    && CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go test -c -o e2e.test ./e2e

# Create tests image
FROM debian:12.5-slim
ARG TARGETOS
ARG TARGETARCH
ARG GIT_TAG
ARG GIT_COMMIT
ENV DD_GIT_TAG=${GIT_TAG}
ENV DD_GIT_COMMIT_SHA=${GIT_COMMIT}
ENV N_VERSION=9.2.0
ENV NODE_VERSION=16.13.0
ENV DATADOG_CI_VERSION=2.32.0
ENV GOTESTSUM_VERSION=1.11.0
ENV GOTESTSUM_ARCHIVE=gotestsum_${GOTESTSUM_VERSION}_${TARGETOS:-linux}_${TARGETARCH}.tar.gz

WORKDIR /

RUN apt-get update \
    && apt-get -y install --no-install-recommends curl ca-certificates \
    && curl --retry 5 -L https://raw.githubusercontent.com/tj/n/v${N_VERSION}/bin/n -o n \
    && bash n ${NODE_VERSION} \
    && npm install -g @datadog/datadog-ci@${DATADOG_CI_VERSION} \
    && curl --retry 5 -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/${GOTESTSUM_ARCHIVE} -o gotestsum.tar.gz \
    && tar -xzf gotestsum.tar.gz \
    && rm -rf ${GOTESTSUM_ARCHIVE}

COPY --from=builder /workspace/e2e.test .
COPY --from=builder /workspace/test2json .
COPY --from=builder /workspace/e2e-runner.sh .
USER 65532:65532

ENTRYPOINT ["/e2e-runner.sh"]
