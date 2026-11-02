# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build single deployable binary
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/bin/redirector ./cmd/api

# Final Stage: Ultra-lightweight runner (<50MB RAM footprint)
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/bin/redirector /app/redirector
COPY --from=builder /app/docs/openapi.json /app/docs/openapi.json
COPY .env.example /app/.env

EXPOSE 8080

ENTRYPOINT ["/app/redirector"]
