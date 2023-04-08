FROM golang:1.20.1

WORKDIR /app
COPY go.mod go.sum *.go ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o ./proxy

CMD ["./proxy"]