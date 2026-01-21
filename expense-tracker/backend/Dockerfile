# Build stage
FROM public.ecr.aws/amazonlinux/amazonlinux:2023 AS builder

# Install Go
RUN yum update -y && yum install -y golang git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Production stage
FROM public.ecr.aws/amazonlinux/amazonlinux:2023

RUN yum update -y && yum install -y ca-certificates tzdata && yum clean all
WORKDIR /root/

# Copy binary and templates
COPY --from=builder /app/main .
COPY --from=builder /app/templates ./templates/

# Expose port
EXPOSE 8081

# Run the application
CMD ["./main"]