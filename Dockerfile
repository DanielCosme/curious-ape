FROM golang:alpine AS base-builder
ENV GOCACHE=/root/.cache/go-build
ENV CGO_ENABLED=0
ARG APE_VERSION=unknown

FROM base-builder AS ape-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
RUN go mod verify
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-extldflags=-static -X main.version=${APE_VERSION}" -o=./bin/ape ./cmd/web

FROM alpine:latest AS ape
RUN apk add --no-cache tzdata
ENV TZ=America/Toronto
# Dedicated group and user with the exact UID/GID
RUN addgroup -g 65532 -S appgroup && \
    adduser -u 65532 -G appgroup -S -D -H appuser && \
    mkdir -p /app && \
    chown -R appuser:appgroup /app

WORKDIR /app
COPY --chown=appuser:appgroup --from=ape-builder /app/bin/ape /app/bin/ape
USER 65532:65532
CMD ["/app/bin/ape"]
