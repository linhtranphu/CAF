# Expense Tracker - Ứng dụng Quản lý Chi phí

Ứng dụng web responsive có thể hoạt động offline, sử dụng AI để phân tích tin nhắn chi phí và lưu trữ vào Google Sheets.

## Tính năng

- ✅ Web app responsive
- ✅ Hoạt động offline (PWA)
- ✅ AI phân tích ngữ nghĩa tin nhắn
- ✅ Tự động lưu vào Google Sheets
- ✅ Đồng bộ khi có mạng

## Cài đặt

### 1. Backend (Go)

```bash
cd backend
go mod tidy
```

### 2. Cấu hình Google Sheets

1. Tạo project trên [Google Cloud Console](https://console.cloud.google.com/)
2. Bật Google Sheets API
3. Tạo Service Account và tải credentials JSON
4. Chia sẻ Google Sheet với email của Service Account
5. Đặt biến môi trường:

```bash
export GOOGLE_CREDENTIALS_FILE=path/to/credentials.json
```

6. Cập nhật `SPREADSHEET_ID` trong `sheets.go`

### 3. Chạy ứng dụng

```bash
# Backend
cd backend
go mod tidy
go run cmd/main.go

# Frontend (serve static files)
cd frontend
# Sử dụng web server bất kỳ, ví dụ:
python -m http.server 3000
```

## Sử dụng

1. Mở http://localhost:3000
2. Nhập tên người dùng
3. Gõ tin nhắn chi phí, ví dụ: "cọc nhà 34 triệu"
4. Ấn "Ghi nhận chi phí"

## Ví dụ tin nhắn

- "cọc nhà 34 triệu" → Items: "cọc nhà", Amount: 34000000
- "mua xăng 200k" → Items: "mua xăng", Amount: 200000
- "ăn trưa 50 nghìn" → Items: "ăn trưa", Amount: 50000

## Cấu trúc dự án (DDD)

```
expense-tracker/
├── backend/
│   ├── cmd/
│   │   └── main.go              # Entry point
│   ├── domain/
│   │   ├── expense/
│   │   │   ├── expense.go       # Expense entity
│   │   │   └── repository.go    # Repository interface
│   │   └── user/
│   │       └── user.go          # User entity
│   ├── application/
│   │   └── services/
│   │       └── expense_service.go # Business logic
│   ├── infrastructure/
│   │   ├── ai/
│   │   │   └── parser.go        # AI message parser
│   │   └── sheets/
│   │       └── repository.go    # Google Sheets impl
│   ├── interfaces/
│   │   └── http/
│   │       ├── expense_handler.go # HTTP handlers
│   │       └── router.go        # Routes
│   └── go.mod
└── frontend/
    ├── index.html
    ├── app.js
    ├── sw.js
    └── manifest.json
```