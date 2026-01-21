# --- GIAI ĐOẠN BUILD (Build stage) ---
# Thay đổi từ 1.23 sang 1.24 để đáp ứng yêu cầu của google.golang.org/genai
FROM golang:1.24-alpine AS builder

# Cài đặt git vì một số thư viện Go (như genai) có thể cần git để tải về
RUN apk add --no-cache git

WORKDIR /app

# Copy các file quản lý thư viện trước để tận dụng Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ mã nguồn
COPY . .

# Build ứng dụng
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go


# --- GIAI ĐOẠN CHẠY (Production stage) ---
FROM public.ecr.aws/amazonlinux/amazonlinux:2023.6.20241121.0

# Cập nhật và cài đặt các gói cần thiết
RUN yum update -y && \
    yum install -y ca-certificates tzdata && \
    yum clean all && \
    rm -rf /var/cache/yum

WORKDIR /root/

# Copy file thực thi từ stage builder
COPY --from=builder /app/main .

# Copy thư mục templates (chỉ thực hiện nếu bạn thực sự có thư mục này)
# Nếu không có templates, hãy xóa hoặc comment dòng dưới đây để tránh lỗi build
COPY --from=builder /app/templates ./templates/

# Mở cổng 8081
EXPOSE 8081

# Lệnh chạy ứng dụng
CMD ["./main"]