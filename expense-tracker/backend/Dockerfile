# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Production stage
FROM public.ecr.aws/amazonlinux/amazonlinux:2023.6.20241121.0

RUN yum update -y && \
    yum install -y ca-certificates tzdata && \
    yum clean all && \
    rm -rf /var/cache/yum

WORKDIR /root/

# Copy binary and templates
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates/

# Expose port
EXPOSE 8081

# Run the application
CMD ["./main"]