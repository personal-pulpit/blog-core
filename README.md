# Blog Core

A robust blog platform rest API written in Go, featuring a clean architecture with comprehensive user management, authentication, and content handling capabilities.

## 🌟 Features

- **Authentication System**
  - JWT-based authentication
  - Redis for session management
  - Secure password hashing
  
- **User Management**
  - User registration and login
  - Email verification system
  - Welcome email functionality
  - User Profile Managements:Delete,Upload,Get
  
- **Article Management**
  - CRUD operations for blog posts
  - Search functionality
  - PostgreSQL for data persistence
  - Likes and Bookmarks by current user

- **Comment Management**
  - Commenting system for blog articles  

- **Category Management**
  - Categories for each article
 
- **Infrastructure**
  - PostgreSQL database
  - Redis cache
  - MinIO Object Storage
  - Email service integration
  - Structured logging with Zap

## 🚀 Getting Started

### Prerequisites
- Go 1.19+
- PostgreSQL
- Redis
- MinIO 
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

## TODO
- Complete User Service tests
- Refactor Swagger  Docs
- Recommendation Articles system
- Include Data for articles
  
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
- [Swagger](https://swagger.io/)
- [Zap](https://github.com/uber-go/zap)
- [MinIO](https://github.com/minio/minio)


For more information or questions, please open an issue in the repository.
