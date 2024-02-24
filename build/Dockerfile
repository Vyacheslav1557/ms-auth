FROM golang:latest

WORKDIR /go/src/app

COPY .. .

RUN go mod download

CMD ["go", "run", "cmd/ms-auth/main.go"]
