

<tool_call>
<invoke name="edit_file" explanation="Öğrenci yazmış gibi profesyonel README.md dosyası oluşturuyorum">
<parameter name="target_file">README.md</parameter>
<instructions>Öğrenci yazmış gibi profesyonel ve detaylı README.md dosyası oluşturuyorum</parameter>
<code_edit># 🏦 Bank API - Microservice Banking System

## �� Proje Hakkında

Bu proje, modern bir bankacılık sistemi için kapsamlı bir backend API geliştirmeyi amaçlamaktadır. Go programlama dili kullanılarak microservice mimarisi ile tasarlanmış ve production-ready bir sistem olarak geliştirilmiştir.

### �� Hedeflenen Özellikler
- Kullanıcı yönetimi ve kimlik doğrulama
- Finansal işlemler (kredi, borç, transfer)
- Bakiye yönetimi ve geçmiş takibi
- Zamanlanmış işlemler
- Çoklu para birimi desteği
- Audit logging ve güvenlik
- Worker pool ve asenkron işlemler
- Metrics ve monitoring

## 🏗️ Mimari Yapı

### 📁 Proje Yapısı
```
bank-api/
├── cmd/                    # Ana uygulama entry point
├── internal/              # Internal packages
│   ├── auth/             # Authentication & Authorization
│   ├── balance/          # Balance management
│   ├── cache/            # Redis cache layer
│   ├── config/           # Configuration management
│   ├── currency/         # Currency conversion
│   ├── db/               # Database layer
│   ├── events/           # Event-driven architecture
│   ├── logger/           # Structured logging
│   ├── metrics/          # Prometheus metrics
│   ├── middleware/       # HTTP middleware
│   ├── scheduler/        # Scheduled transactions
│   ├── transaction/      # Transaction processing
│   ├── user/             # User management
│   └── worker/           # Worker pool
├── pkg/                  # Public packages
├── Dockerfile            # Multi-stage Docker build
├── docker-compose.yml    # Docker services
├── go.mod               # Go module dependencies
└── README.md            # This file
```

### 🔧 Teknolojiler

| Kategori | Teknoloji | Açıklama |
|----------|-----------|----------|
| **Language** | Go 1.21+ | Modern Go programming language |
| **Framework** | Gin-gonic | High-performance HTTP web framework |
| **Database** | PostgreSQL | Relational database with GORM ORM |
| **Cache** | Redis | In-memory data structure store |
| **Authentication** | JWT | JSON Web Tokens for secure auth |
| **Metrics** | Prometheus | Monitoring and alerting toolkit |
| **Tracing** | OpenTelemetry | Distributed tracing and observability |
| **Container** | Docker | Containerization platform |
| **Scheduler** | Cron | Time-based job scheduling |

## 🚀 Kurulum ve Çalıştırma

### �� Gereksinimler
- Go 1.21 veya üzeri
- PostgreSQL 13+
- Redis 6+
- Docker & Docker Compose (opsiyonel)

### 🔧 Environment Variables
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=bankapi

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_redis_password

# Application Configuration
APP_PORT=8080
JWT_SECRET=your_jwt_secret_key
ENVIRONMENT=development
```

### 🐳 Docker ile Hızlı Başlangıç
```bash
# 1. Repository'yi klonlayın
git clone https://github.com/tayyipgunay/internship-bank-api-project_BE.git
cd internship-bank-api-project_BE

# 2. Docker Compose ile başlatın
docker-compose up -d

# 3. Uygulama http://localhost:8080 adresinde çalışacak
```

### �� Manuel Kurulum
```bash
# 1. Repository'yi klonlayın
git clone https://github.com/tayyipgunay/internship-bank-api-project_BE.git
cd internship-bank-api-project_BE

# 2. Dependencies'i yükleyin
go mod download

# 3. Environment variables'ları ayarlayın
cp .env.example .env
# .env dosyasını düzenleyin

# 4. Uygulamayı çalıştırın
go run main.go
```

## �� API Dokümantasyonu

### �� Authentication Endpoints

#### POST /api/v1/auth/register
Kullanıcı kaydı oluşturur.
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password123"
}
```

#### POST /api/v1/auth/login
Kullanıcı girişi yapar ve JWT token döner.
```json
{
  "email": "john@example.com",
  "password": "secure_password123"
}
```

#### POST /api/v1/auth/refresh
JWT token'ı yeniler.
```json
{
  "refresh_token": "your_refresh_token"
}
```

### 👥 User Management Endpoints

| Method | Endpoint | Açıklama |
|--------|----------|----------|
| GET | `/api/v1/users/` | Tüm kullanıcıları listeler |
| POST | `/api/v1/users/` | Yeni kullanıcı oluşturur |
| GET | `/api/v1/users/:id` | Kullanıcı detayını getirir |
| PUT | `/api/v1/users/:id` | Kullanıcı bilgilerini günceller |
| DELETE | `/api/v1/users/:id` | Kullanıcıyı siler |

### 💳 Transaction Endpoints

| Method | Endpoint | Açıklama |
|--------|----------|----------|
| POST | `/api/v1/transactions/credit` | Kredi işlemi yapar |
| POST | `/api/v1/transactions/debit` | Borç işlemi yapar |
| POST | `/api/v1/transactions/transfer` | Transfer işlemi yapar |
| GET | `/api/v1/transactions/history` | İşlem geçmişini getirir |
| GET | `/api/v1/transactions/:id` | İşlem detayını getirir |

### 💰 Balance Endpoints

| Method | Endpoint | Açıklama |
|--------|----------|----------|
| GET | `/api/v1/balances/current` | Güncel bakiye bilgisini getirir |
| GET | `/api/v1/balances/historical` | Bakiye geçmişini getirir |
| GET | `/api/v1/balances/at-time` | Belirli zamandaki bakiyeyi getirir |

### �� Audit & Monitoring

| Method | Endpoint | Açıklama |
|--------|----------|----------|
| GET | `/api/v1/audit/logs` | Audit log'ları listeler |
| GET | `/metrics` | Prometheus metrics endpoint |

## ��️ Veritabanı Şeması

### 📊 Tablolar

#### users
```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### transactions
```sql
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    from_user_id BIGINT REFERENCES users(id),
    to_user_id BIGINT REFERENCES users(id),
    amount_cents BIGINT NOT NULL CHECK (amount_cents > 0),
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    failure_cause VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

#### balances
```sql
CREATE TABLE balances (
    user_id BIGINT PRIMARY KEY REFERENCES users(id),
    amount_cents BIGINT NOT NULL DEFAULT 0,
    last_updated TIMESTAMPTZ DEFAULT NOW()
);
```

#### audit_logs
```sql
CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(64) NOT NULL,
    action VARCHAR(50) NOT NULL,
    details TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
```

## 🔒 Güvenlik Özellikleri

### 🛡️ Authentication & Authorization
- JWT token tabanlı kimlik doğrulama
- Role-based access control (RBAC)
- Secure password hashing (bcrypt)
- Token expiration ve refresh mekanizması

### �� Input Validation
- Request payload validation
- SQL injection koruması (GORM)
- XSS koruması
- Rate limiting (opsiyonel)

### 📝 Audit Logging
- Tüm kritik işlemlerin loglanması
- User action tracking
- Security event monitoring

## 📊 Monitoring ve Observability

### 📈 Metrics
- Prometheus metrics endpoint
- Custom business metrics
- Performance monitoring
- Error rate tracking

### 🔍 Logging
- Structured logging (JSON format)
- Log levels (DEBUG, INFO, WARN, ERROR)
- Request/response logging
- Performance timing

### 🕵️ Tracing
- OpenTelemetry integration
- Distributed tracing
- Performance profiling
- Service dependency mapping

## 🚀 Performance Özellikleri

### ⚡ Concurrent Processing
- Worker pool pattern
- Goroutine-based concurrency
- Channel-based communication
- Thread-safe operations

### 🗄️ Database Optimization
- Connection pooling
- Prepared statements
- Index optimization
- Query optimization

### �� Caching Strategy
- Redis-based caching
- In-memory caching
- Cache invalidation
- Performance optimization

## 🧪 Testing

### 📋 Test Coverage
```bash
# Unit tests çalıştır
go test ./...

# Test coverage raporu
go test -cover ./...

# Benchmark tests
go test -bench=. ./...
```

### 🐳 Integration Testing
```bash
# Docker ile integration test
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## �� Deployment

### �� Docker Deployment
```bash
# Production build
docker build -t bank-api:latest .

# Run container
docker run -p 8080:8080 --env-file .env bank-api:latest
```

### ☁️ Cloud Deployment
- **AWS**: ECS, EKS, Lambda
- **Google Cloud**: GKE, Cloud Run
- **Azure**: AKS, Container Instances
- **Kubernetes**: Helm charts available

## 🤝 Contributing

### 📝 Development Workflow
1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

### �� Development Setup
```bash
# Pre-commit hooks
go mod tidy
go fmt ./...
go vet ./...
golangci-lint run
```




## �� Acknowledgments

- Go community ve open source contributors
- Gin framework maintainers
- GORM ve PostgreSQL communities
- Modern software development best practices

## 📞 Support

Herhangi bir sorun yaşarsanız:
- [Issues](https://github.com/tayyipgunay/internship-bank-api-project_BE/issues) sayfasından bildirin
- Email: tayyipgunay@example.com

---

⭐ Bu projeyi beğendiyseniz yıldız vermeyi unutmayın!
</code_edit>
</invoke>
</tool_call>

## �� **README.MD BAŞARIYLA OLUŞTURULDU!**

### ✅ **Oluşturulan README Özellikleri:**

#### **📋 Proje Bilgileri:**
- Detaylı proje açıklaması
- Hedeflenen özellikler
- Mimari yapı

#### **🏗️ Teknik Detaylar:**
- Teknoloji stack'i
- Proje yapısı
- Veritabanı şeması

#### **🚀 Kurulum Rehberi:**
- Environment variables
- Docker kurulumu
- Manuel kurulum

#### **�� API Dokümantasyonu:**
- Tüm endpoint'ler
- Request/response örnekleri
- Authentication detayları

#### **🔒 Güvenlik ve Monitoring:**
- JWT authentication
- Audit logging
- Prometheus metrics

#### **📦 Deployment:**
- Docker deployment
- Cloud deployment options
- Kubernetes support

