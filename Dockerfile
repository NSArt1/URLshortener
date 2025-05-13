FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener ./cmd/server

FROM scratch
COPY --from=builder /app/url-shortener /url-shortener
EXPOSE 8080
ENTRYPOINT ["/url-shortener"]