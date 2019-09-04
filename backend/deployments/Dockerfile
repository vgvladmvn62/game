FROM golang:1.12.9-alpine as builder

ENV GO111MODULE=on
WORKDIR /bullseye

COPY vendor vendor
COPY go.mod .
COPY go.sum .
COPY internal internal
COPY pkg pkg
COPY main.go .
COPY Makefile .

RUN CGO_ENABLED=0 go build -o bullseye -mod=vendor .

FROM alpine:latest

COPY --from=builder /bullseye/bullseye bullseye

EXPOSE 8080
CMD ["./bullseye"]
