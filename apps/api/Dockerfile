# Stage 1: Build
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache gcc musl-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o app-server main.go
RUN CGO_ENABLED=1 go build -o app-migrate cmd/migrate/main.go

# Stage 2: Run
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata curl libstdc++
WORKDIR /app
COPY --from=builder /app/app-server .
COPY --from=builder /app/app-migrate .
COPY --from=builder /app/lang ./lang
COPY --from=builder /app/database ./database

EXPOSE 3001
CMD ["./app-server"]
