FROM golang:1.20.2-alpine3.17

WORKDIR /app

ADD . /app

RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

ENTRYPOINT ["sh", "-c"]
