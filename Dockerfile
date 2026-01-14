FROM golang:1.25.1 AS base-builder
ENV GOCACHE=/root/.cache/go-build
ENV CGO_ENABLED=1
ARG APE_VERSION=unknown

FROM base-builder AS echo-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
RUN go mod verify
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-s -extldflags=-static -X main.version=1.0" -o=./bin/echo ./cmd/echo

FROM alpine:latest AS echo
RUN apk add --no-cache tzdata
RUN apk add gcompat
ENV TZ=America/Toronto
WORKDIR /app
COPY --from=echo-builder /app/bin/echo /app/bin/echo
CMD ["/app/bin/echo"]

FROM base-builder AS ape-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
RUN go mod verify
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-s -extldflags=-static -X main.version=${APE_VERSION}" -o=./bin/ape ./cmd/web

# Main App
FROM alpine:latest AS ape
RUN apk add --no-cache tzdata
RUN apk add gcompat
ENV TZ=America/Toronto
WORKDIR /app
COPY --from=ape-builder /app/bin/ape /app/bin/ape
CMD ["/app/bin/ape"]


FROM base-builder AS migrate-builder
WORKDIR /app
RUN git clone --branch v4.17.0 --depth 1 https://github.com/golang-migrate/migrate
WORKDIR /app/migrate
RUN --mount=type=cache,target=/go/pkg/mod go mod download 
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o build/migrate -ldflags="-s -w -extldflags=-static" -tags 'sqlite3' ./cmd/migrate


# Migrate Container
FROM alpine:latest AS migrate
COPY ./database/migrations/sqlite /migrations
COPY --from=migrate-builder /app/migrate/build/migrate /usr/local/bin/migrate
ENTRYPOINT ["migrate"]
CMD ["--help"]
