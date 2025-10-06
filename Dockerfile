# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w -extldflags "-static"' -o app main.go

# Final stage
FROM alpine:latest

# RUN addgroup -S nonroot \
#     && adduser -S nonroot -G nonroot
# USER nonroot

WORKDIR /app

# Copy compiled binary from builder
COPY --from=builder /app/app /app/app

RUN apk --no-cache add tzdata ca-certificates
ENV TZ=Asia/Jakarta

EXPOSE 8080
ENTRYPOINT ["./app"]
