FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/client
RUN go build -o client main.go

CMD ["./client", "weighted"]