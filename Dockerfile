FROM golang:alpine as builder

WORKDIR /go/src/app
COPY *.go .

RUN apk add --no-cache git
RUN go get -d -t ./...
RUN go build

FROM alpine

WORKDIR /var/app
COPY --from=builder /go/src/app /var/app/

CMD "./app"
