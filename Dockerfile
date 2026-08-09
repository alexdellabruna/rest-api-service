FROM golang:alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

FROM alpine:3.23

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1000 appuser && adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/main .

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./main"]
