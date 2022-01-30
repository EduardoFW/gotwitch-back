FROM golang:latest

WORKDIR /go/app

COPY . /go/app

RUN go build -o /go/app/main

CMD ["/go/app/main"]