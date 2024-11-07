FROM golang:1.23.3-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go build -o /mywalletapp ./cmd

RUN apk --no-cache add ca-certificates bash
RUN apk add --no-cache --virtual .build-deps gcc libc-dev
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

EXPOSE 8080
