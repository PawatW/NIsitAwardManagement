# ---------- build stage ----------
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache ca-certificates
RUN go build -o server ./cmd/api

# ---------- runtime stage ----------
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/server /app/server
EXPOSE 8080
CMD ["/app/server"]