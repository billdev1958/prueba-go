FROM golang:1.26.1-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main cmd/main.go

FROM alpine:3.23.3

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/db/schema.sql ./db/schema.sql

ENV PORT=8080

EXPOSE 8080

CMD ["./main"]
