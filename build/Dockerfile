FROM golang:alpine3.15 as builder

ARG APPLICATION=ui-assignment

RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base && apk add docker

COPY . /go/${APPLICATION}

WORKDIR /go/${APPLICATION}

RUN go get ${APPLICATION}

RUN go install

RUN go build -o ./bin/ui-assignment ./cmd/server/...

EXPOSE 8080

