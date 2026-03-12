FROM golang:1.26 AS builder
WORKDIR /app
COPY . .
RUN go build -v -o server .

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/server .
CMD ["./server"]