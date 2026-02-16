# -------- STAGE 1: build --------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/main.go

# -------- STAGE 2: runtime --------
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8083

CMD ["./app"]
