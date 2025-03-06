# Blog Core

A robust blog platform backend written in Go, featuring a clean architecture with comprehensive user management, authentication, and content handling capabilities.

## 🌟 Features

- **Authentication System**
  - JWT-based authentication
  - Redis for session management
  - Secure password hashing
  
- **User Management**
  - User registration and login
  - Email verification system
  - Welcome email functionality
  
- **Article Management**
  - CRUD operations for blog posts
  - PostgreSQL for data persistence
  
- **Infrastructure**
  - PostgreSQL database
  - Redis cache
  - Email service integration
  - Structured logging with Zap

## 🚀 Getting Started

### Prerequisites
- Go 1.19+
- PostgreSQL
- Redis
- Make (optional)

### Installation

```bash
git clone https://github.com/personal-pulpit/blog-core.git
cd blog-core
go mod download
```

### Running the Application

```bash
go run main.go
```

## 🏗️ Project Structure

```
blog-core/
├── api/
│   └── server/        # HTTP server and routing
├── config/           # Configuration management
├── database/
│   ├── postgres/     # PostgreSQL connection and repositories
│   └── redis/        # Redis connection and operations
├── internal/
│   └── service/      # Business logic services
├── pkg/
│   ├── auth_manager/ # Authentication utilities
│   ├── email_manager/# Email service
│   └── logger/       # Logging utilities
└── utils/            # Common utilities
```

## 🔄 API Endpoints

- **Authentication**
  - POST `/auth/register`
  - POST `/auth/login`
  - POST `/auth/verify-email`

- **Users**
  - GET `/users/profile`
  - PUT `/users/update`

- **Articles**
  - GET `/articles`
  - POST `/articles`
  - PUT `/articles/{id}`
  - DELETE `/articles/{id}`

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch:
```bash
git checkout -b feature/amazing-feature
```

3. Commit your changes:
```bash
git commit -m "Add amazing feature"
```

4. Push to the branch:
```bash
git push origin feature/amazing-feature
```

5. Open a Pull Request

## 📝 Development Guidelines

- Follow Go best practices and idioms
- Write tests for new features
- Update documentation when adding features
- Use meaningful commit messages
- Follow the existing code structure

## 🛠️ Built With

- [Go](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [JWT](https://jwt.io/)
- [Zap](https://github.com/uber-go/zap)

## 🙏 Acknowledgments

- Thanks to all contributors
- Inspired by clean architecture principles
- Built with modern Go practices

For more information or questions, please open an issue in the repository.
