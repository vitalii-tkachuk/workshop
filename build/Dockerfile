FROM golang:1.19-alpine

RUN apk update \
   && apk add --no-cache git \
   && go install github.com/codegangsta/gin@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

ENTRYPOINT ["gin", "--appPort=8080", "--port=80", "run"]
