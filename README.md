# WarOnk Backend (Go + gRPC + REST)

Project ini menggunakan **Clean Architecture** dengan implementasi **REST API** dan **gRPC**.

---

## 📂 Struktur Folder

```
internal/
│── delivery/
│   ├── grpc/            # Handler gRPC
│   └── presenter/       # Response DTO / presenter
│
│── domain/
│   ├── entity/          # Struct entity domain
│   ├── repository/      # Interface repository
│   └── usecase/         # Business logic (use cases)
│
│── repository/
│   ├── db/              # Implementasi repository (DB)
│   └── transaction/     # Manajemen transaksi DB
│
│── infrastructure/
│   ├── config/          # Konfigurasi project (env, dsb)
│   ├── database/        # Inisialisasi koneksi database
│
│── proto/               # File .proto untuk gRPC
│── cmd/                 # Entry point aplikasi
│   └── main.go
```

---

## 🚀 Cara Menjalankan

### 1. Clone Repository

```bash
git clone https://github.com/Temisaputra/warOnk.git
cd warOnk
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Setup Database

Buat database baru, lalu sesuaikan konfigurasi di file:

```
internal/infrastructure/config/config.go
```

### 4. Generate gRPC Code

Jika ada perubahan di file `.proto`:

```bash
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

### 5. Run gRPC Server

```bash
go run cmd/main.go grpc
```

---

## 📌 Catatan

- `internal/domain` adalah **core business logic**.
- `internal/repository` adalah implementasi akses data.
- `internal/delivery` adalah interface untuk komunikasi (gRPC).
- `internal/infrastructure` adalah konfigurasi teknis.

---

## 🛠 Tools

- **Go 1.22+**
- **gRPC**
- **PostgreSQL**
- **GORM**
- **Protobuf Compiler**

# gRPC + Envoy Setup

## Cara Menjalankan

1. Build dan jalankan Docker Compose:

   ```bash
   docker-compose up --build
   ```

2. Service yang tersedia:
   - gRPC Service → `localhost:50051`
   - Envoy Proxy (gRPC-Web) → `localhost:8081`
   - Envoy Admin → `localhost:9901`

## Testing dari Client

### 1. Install grpcurl

Untuk mengetes gRPC dari terminal:

```bash
# MacOS / Linux (Homebrew)
brew install grpcurl

# Ubuntu/Debian
sudo apt-get install grpcurl
```

### 2. Testing gRPC langsung ke service

```bash
grpcurl -plaintext localhost:50051 list
```

### 3. Testing gRPC melalui Envoy (gRPC-Web)

```bash
grpcurl -plaintext -proto proto/user.proto -d '{"id": "123"}' localhost:8081 user.UserService/GetUser
```

### 4. Testing via Web Browser (opsional)

Bisa gunakan gRPC-Web client di React/Vue dengan konfigurasi endpoint ke `http://localhost:8081`.

---

## gRPC Testing

### 1. Using grpcurl (without TLS)

### List grpc service

```bash
grpcurl -plaintext localhost:50051 list
```

### GET ALL

```bash
grpcurl -plaintext -d '{"id": "1"}' localhost:50051 user.UserService/GetUser
```

### GET BY ID

```bash
grpcurl -plaintext -d '{"id": "7"}' localhost:50051 userpb.UserService/GetUserById
```

### CREATE USER

```bash
grpcurl -plaintext -d '{
  "name": "Hamdi Created",
  "email": "hamdi.created@gmail.com",
  "role": "superadmin"
}' localhost:50051 userpb.UserService/CreateUser
```

### UPDATE USER

```bash
grpcurl -plaintext -d '{
  "id": 7,
  "name": "Rahman Updated",
  "email": "rahman.updated@gmail.com",
  "role": "superadmin"
}' \
localhost:50051 userpb.UserService/UpdateUser
```

### DELETE USER

```bash
grpcurl -plaintext -d '{"id": 7}' localhost:50051 userpb.UserService/DeleteUser
```

### 2. Example gRPC Request (GetUser)

```bash
grpcurl -plaintext -d '{"id": "123"}' localhost:50051 user.UserService/GetUser
```

Expected response:

```json
{
  "id": "123",
  "name": "Temi Saputra",
  "email": "temi@example.com"
}
```

### Generate gRPC Client (Frontend)

Jalankan perintah berikut untuk generate client di sisi frontend (gRPC-Web):

```bash
protoc -I=./proto   proto/user.proto   --js_out=import_style=commonjs:./frontend   --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./frontend
```

Selamat mencoba 🚀
