FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/service
RUN go build -o service main.go

CMD ["./service", "S1", "8001", "3"]