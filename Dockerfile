FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /moodmap-api ./cmd/api

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /moodmap-api /app/moodmap-api

EXPOSE 8080

CMD ["./moodmap-api"]
