FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk update && apk add --no-cache git gcc musl-dev
COPY producer/go.* .
RUN go mod download && go mod verify
COPY producer/. .
WORKDIR /app
RUN CGO_ENABLED=1 go build -tags musl -o producer .

FROM alpine:3.20
COPY --from=builder /app/producer /usr/bin/producer
CMD ["/usr/bin/producer"]
