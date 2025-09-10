# Stage 1: Builder
FROM golang:1.25.1-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/server/

# Stage 2: Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/frontend ./frontend
COPY config.yml .

EXPOSE 8080

CMD ["./main"]