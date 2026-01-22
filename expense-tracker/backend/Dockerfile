# Build stage
FROM golang:1.24-alpine AS builder

# Cài đặt git vì một số thư viện Go (như genai) có thể cần git để tải về
RUN apk add --no-cache git

WORKDIR /app

# Copy các file quản lý thư viện trước để tận dụng Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ mã nguồn
COPY . .

# Cập nhật go.mod để tương thích với Go 1.24
RUN go mod tidy

# Build ứng dụng
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