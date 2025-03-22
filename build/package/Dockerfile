FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o bot ./cmd/archive_bot

FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /root

COPY --from=builder /app/bot .

COPY configs/ ./configs/

CMD [ "./bot", "--config", "./configs/prod.yaml"]