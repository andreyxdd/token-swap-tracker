FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk update && apk add --no-cache git gcc musl-dev
COPY consumer/go.* .
RUN go mod download && go mod verify
COPY consumer/. .
WORKDIR /app/cmd/api
RUN CGO_ENABLED=1 go build -tags musl -o api .

FROM alpine:3.20
RUN apk add -U --no-cache ca-certificates
COPY --from=builder /app/cmd/api/api /usr/bin/api
ENV PORT=8081
EXPOSE ${PORT}
CMD ["/usr/bin/api"]
