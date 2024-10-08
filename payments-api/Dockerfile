FROM golang:1.23.2-alpine AS builder
WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./

COPY main.go ./
COPY main_test.go ./

RUN go mod download

RUN go build ./main.go
RUN chmod a+x main

FROM alpine:latest
WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/main /usr/src/app/main

CMD ["./main"]