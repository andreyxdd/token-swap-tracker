FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk update && apk add --no-cache git gcc musl-dev
COPY producer/go.* .
RUN go mod download && go mod verify
COPY producer/. .
WORKDIR /app/cmd
RUN CGO_ENABLED=1 go build -tags musl -o producer .
RUN ls -al 
RUN ls -al /app

FROM alpine:3.20
ENV KAFKA_BROKERS=localhost:9092
ENV KAFKA_TOPIC=swaps
COPY --from=builder /app/cmd/producer /usr/bin/producer
CMD ["/usr/bin/producer"]
