FROM golang:alpine AS base-builder
ENV GOCACHE=/root/.cache/go-build
ENV CGO_ENABLED=0
ARG APE_VERSION=unknown

FROM base-builder as ape-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
RUN go mod verify
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-s -w -extldflags=-static -X main.version=${APE_VERSION}" -o=./bin/ape ./cmd/web

FROM alpine:latest AS ape
RUN apk add --no-cache tzdata
ENV TZ=America/Toronto
WORKDIR /app
COPY --from=ape-builder /app/bin/ape /app/bin/ape
CMD ["/app/bin/ape"]
