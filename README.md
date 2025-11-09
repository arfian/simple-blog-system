# simple-blog-system
simple-blog-system

A modern application API build with Go, featuring attendance, overtime, reimbusment and payslip

## Prerequisites

Before running the application, make sure you have the following installed:

- Go 1.21 or higher
- PostgreSQL
- Golang Migrate

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/arfian/simple-blog-system.git
cd simple-blog-system
```

### 2. Environment Setup

Create a `.env` file in the root directory, you can copy value .env.example

```bash
APP_NAME=simple-blog-system
APP_ENV=local
APP_PORT=8089

DB_DSN=
DB_MAX_OPEN_CONN=100
DB_MAX_IDLE_CONN=10
DB_MAX_LIFETIME_CONN=4
DB_MAX_IDLETIME_CONN=1

SIGNING_KEY=simpleblogsystem123
CACHE_TTL=10
```

### 3. Go-migrate CLI
```sh
#mac
$ brew install golang-migrate

#linux
$ curl -L https://github.com/golang-migrate/migrate/releases/download/$version/migrate.$platform-amd64.tar.gz | tar xvz
```

### 4. Running Migration
#### Migration Up
```sh
$ make migrateup
```

#### Migration Down
```sh
$ make migratedown
```

## How To Run
### Using Makefile
```sh
$ make run 
```

### Using Terminal / Cmd
```sh
$ go mod download
$ go run main.go 
```

### Check Unit Test
```sh
$ make test 
```

## ERD Database
ERD Database you can click url dbdiagram : https://dbdiagram.io/d/simple-blog-system-690fbf026735e11170e11198

## API Documentation
The API documentation is available in Postman format. Import the following files into Postman:

- `postman/Simple Blog Post.postman_collection.json`

### Key Endpoints

1. User
   - POST `/v1/public-api/user/register` - User registration
   - POST `/v1/public-api/user/login` - User login
   - GET `/v1/api/profile/` - User login

2. Post
   - POST `/v1/api/post` - insert post data
   - PUT `/v1/api/post/{id}` - update post data
   - DELETE `/v1/api/post/{id}` - Delete post data
   - GET `/v1/api/post/{id}` - Get Post By ID
   - GET `/v1/api/post` - Get All Post

3. Comment
   - POST `/v1/api/comment` - insert comment data
   - PUT `/v1/api/comment/{id}` - update comment data
   - DELETE `/v1/api/comment/{id}` - Delete comment data
   - GET `/v1/api/comment/{id}` - Get Pcommentost By ID
   - GET `/v1/api/comment` - Get All comment

## Project Structure
```
.
├── cmd/
│   ├── rest/                # Setup rest API with GIN
│   ├── ├── middleware/      # Register list API url by domain
│   ├── pubsub/              # Setup pubsub
├── config/                  # Setup config register env application
│   ├── db/                  # Config init database
├── docs/                    # Generate docs API swagger
├── external/                # Code with external system
├── internal/                # Code with internal system
│   ├── app/                 # Bussiness domain
│   ├── ├── {domain_name}/   # Folder domain name
│   ├── ├── ├── handler/     # Logic handle API logic
│   ├── ├── ├── model/       # Model struct of domain
│   ├── ├── ├── payload/     # Model struct payload response or param
│   ├── ├── ├── port/        # List all interface domain
│   ├── ├── ├── repository/  # Logic query SQL
│   ├── ├── ├── server/      # Register endpoint API with spesific domain
│   ├── ├── ├── service/     # Business logic implementation
│   ├── setup/               # Setup init register domain interfaces
└── log/                     # Generate log file
└── migrations/              # Generate code migration database
└── pkg/                     # code package helper logic
└── postman/                 # Postman collection and environment
```

## Technologies
- [Golang](https://go.dev/)
- [Gorm](https://gorm.io/index.html)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [Swaggo](https://github.com/swaggo/swag)
- [Gin](https://gin-gonic.com/)
- [Zerolog](https://github.com/rs/zerolog)
- PostgreSQL

## Accessing Swagger
```
localhost:8089/swagger/index.html
```