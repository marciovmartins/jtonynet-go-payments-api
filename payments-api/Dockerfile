FROM golang:1.23.2-alpine AS builder
WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./

COPY . .

RUN go mod download

RUN go build ./cmd/http/main.go
RUN chmod a+x main

FROM alpine:latest
WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/.env /usr/src/app/.env
COPY --from=builder /usr/src/app/main /usr/src/app/main

CMD ["./main"]