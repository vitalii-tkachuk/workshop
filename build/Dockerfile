FROM golang:1.19-alpine

RUN go install github.com/codegangsta/gin@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ENTRYPOINT ["gin", "--port", "80", "--appPort", "8080", "--build", "cmd/server", "--path", "."]
