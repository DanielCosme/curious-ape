FROM golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download 
RUN go mod verify
COPY . .
ENV CGO_ENABLED=1
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build -ldflags="-extldflags=-static" -o=./bin/ape ./cmd/web


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/config.json /app/config.json
COPY --from=builder /app/bin/ape /app/bin/ape
CMD [ "/app/bin/ape", "-env", "prod" ]