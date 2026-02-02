FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum* ./

COPY main.go ./

RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service .

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

WORKDIR /root/

COPY --from=builder /app/service .

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./service"]
