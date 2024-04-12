FROM golang:1.22 as base-builder
ENV GOCACHE=/root/.cache/go-build
ENV CGO_ENABLED=1
ARG APE_VERSION=unknown

# Ape

FROM base-builder as ape-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download 
RUN go mod verify
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-s -extldflags=-static -X main.version=${APE_VERSION}" -o=./bin/ape ./cmd/web


FROM alpine:latest as ape
WORKDIR /app
COPY --from=ape-builder /app/bin/ape /app/bin/ape
CMD ["/app/bin/ape"]

# Migrate

FROM base-builder as migrate-builder
WORKDIR /app
RUN git clone --branch v4.17.0 --depth 1 https://github.com/golang-migrate/migrate
WORKDIR /app/migrate
RUN --mount=type=cache,target=/go/pkg/mod go mod download 
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -o build/migrate -ldflags="-s -w -extldflags=-static" -tags 'sqlite3' ./cmd/migrate


FROM alpine:latest AS migrate
COPY ./migrations/sqlite /migrations
COPY --from=migrate-builder /app/migrate/build/migrate /usr/local/bin/migrate
ENTRYPOINT ["migrate"]
CMD ["--help"]
