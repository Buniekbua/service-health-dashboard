FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o service-health-dashboard ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/service-health-dashboard .

EXPOSE 8080

CMD ["./service-health-dashboard"]