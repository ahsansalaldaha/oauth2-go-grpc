FROM golang:1.18-alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

COPY go.mod ./
COPY go.sum ./
RUN go mod download

CMD ["air", "-c", ".air.toml"]
