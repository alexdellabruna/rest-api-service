FROM golang:alpine3.23 AS builder

RUN apk add --no-cache build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o task-red-hat ./cmd

FROM alpine:3.23

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/task-red-hat .

RUN chown -R appuser:appuser /app

USER appuser

CMD ["./task-red-hat"]
