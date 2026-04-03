FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o sweepline ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/sweepline .

EXPOSE 5200
CMD ["./sweepline"]
