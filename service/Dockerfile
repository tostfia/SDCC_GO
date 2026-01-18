FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

# Crea la directory per lo stato condiviso
RUN mkdir -p /app/state

WORKDIR /app/service
RUN go build -o service main.go

# Il nome del servizio viene da SERVICE_NAME nel docker-compose
CMD ["./service", "8001", "3"]