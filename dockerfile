FROM golang:1.26.1-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o guild_tracker .

FROM alpine:3.21

RUN apk add --no-cache -ca-certificates

WORKDIR /app

COPY --from=builder /app/guild_tracker .

VOLUME [ "/app/database" ]

CMD [ "./guild_tracker" ]
