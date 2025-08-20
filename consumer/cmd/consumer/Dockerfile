FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk update && apk add --no-cache git gcc musl-dev
COPY consumer/go.* .
RUN go mod download && go mod verify
COPY consumer/. .
WORKDIR /app/cmd/consumer
RUN CGO_ENABLED=1 go build -tags musl -o consumer .

FROM alpine:3.20
RUN apk add -U --no-cache ca-certificates
COPY --from=builder /app/cmd/consumer/consumer /usr/bin/consumer
CMD ["/usr/bin/consumer"]
