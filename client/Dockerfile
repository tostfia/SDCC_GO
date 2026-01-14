FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

WORKDIR /app/client
RUN go build -o client .

CMD ["./client", "weighted"]