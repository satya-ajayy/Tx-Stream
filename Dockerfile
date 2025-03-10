# base image
FROM golang:1.23.3-alpine as base
WORKDIR /go-kafka

ENV CGO_ENABLED=0

COPY go.mod go.sum /go-kafka/
RUN go mod download

ADD . .
RUN go build -o /usr/local/bin/go-kafka ./cmd/go-kafka

# runner image
FROM gcr.io/distroless/static:latest
WORKDIR /app
COPY --from=base /usr/local/bin/go-kafka go-kafka

EXPOSE 4202
ENTRYPOINT ["/app/go-kafka"]
