FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . /app
RUN go generate ./...
RUN go build -o sku /app/cmd/sku/main.go

FROM alpine:3.15
COPY --from=builder /app/sku /
CMD ["/sku"]
