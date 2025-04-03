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
  - User Profile Managements:Delete,Upload,Get
  
- **Article Management**
  - CRUD operations for blog posts
  - Search functionality
  - PostgreSQL for data persistence

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

## 🏗️ Project Structure(shorted)

```
blog-core/
├── api/
│   ├── handlers/   # HTTP request handlers
│   ├── middleware/ # Middleware functions
│   └── routes/     # HTTP routes
│   └── server/        # HTTP server and routing
├── config/           # Configuration management
├── database/
│   ├── postgres/     # PostgreSQL connection and repositories
│   └── redis/        # Redis connection and operations
├── internal/
│   ├── model/        # Data models
│   ├── repository/   # Database repositories
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
  - GET `/auth/authenticate`
  - POST `/auth/refresh-token`
  - POST `/auth/logout`
  - POST `/auth/change-password`
  - POST `/auth//reset-password/request`
  - POST `/auth/reset-password/submit`

- **Users**
  - GET `/users/me`
  - GET `/users/{id}`
  - PATCH `/users/update`
  - DELETE `/users/delete`
  - POST `/users/add-profile-picture` 
  - GET `/users/profile-picture-url/{id}`
  - DELETE `/users/delete-profile-picture`

- **Articles**
  - GET `/articles`
  - POST `/articles/create`
  - GET `/articles/{id}`
  - GET `/articles/search`
  - PATCH `/articles/{id}`
  - DELETE `/articles/{id}`

- **Comments**
  - POST  `/comments/create`
  - PATCH `/comments/{id}`
  - DELETE `/comments/{id}`

- **Categories**
  - POST  `/categories/create`
  - GET  `/categories/` 
  - GET `/categories/{id}`
  - DELETE `/categories/{id}`
   
## 📝 Documentation
After ran go run main.go open: http://localhost:8000/swagger/index.html in your browser

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
- Search Categories by titles
- Bookmarks management
- Likes managemnet
- Tags management  

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

## 🙏 Acknowledgments

- Thanks to all contributors
- Inspired by clean architecture principles
- Built with modern Go practices

For more information or questions, please open an issue in the repository.
