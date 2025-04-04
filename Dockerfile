# base image
FROM golang:1.23.3-alpine as base
WORKDIR /tx-stream

ENV CGO_ENABLED=0
COPY go.mod go.sum /tx-stream/
RUN go mod download

ADD . .
RUN go build -o /usr/local/bin/tx-stream ./cmd/tx-stream

# runner image
FROM gcr.io/distroless/static:latest
WORKDIR /app
COPY --from=base /usr/local/bin/tx-stream tx-stream
ENTRYPOINT ["/app/tx-stream"]
