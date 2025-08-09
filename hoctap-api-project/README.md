# HocTap API

A modern REST API built with Go, Gorilla Mux router, and MySQL database.

## Features

- 🚀 RESTful API endpoints
- 👥 User management (CRUD operations) 
- 💾 MySQL database integration
- 🔐 Environment-based configuration
- 🏥 Health check endpoint with database status
- 🔧 CORS support
- 📝 Request logging middleware
- 📊 JSON responses with timestamps
- 🌐 Integrated HTML dashboard
- 📈 User statistics endpoint

## Getting Started

### Prerequisites

- Go 1.19 or later
- MySQL 5.7+ or MariaDB 10.3+
- Git (optional)

### Database Setup

1. **Install MySQL/MariaDB**
   ```bash
   # On Windows (using Chocolatey)
   choco install mysql

   # On macOS (using Homebrew)
   brew install mysql

   # On Ubuntu/Debian
   sudo apt-get install mysql-server
   ```

2. **Create Database**
   ```sql
   CREATE DATABASE hoctap_api;
   CREATE USER 'hoctap_user'@'localhost' IDENTIFIED BY 'your_secure_password';
   GRANT ALL PRIVILEGES ON hoctap_api.* TO 'hoctap_user'@'localhost';
   FLUSH PRIVILEGES;
   ```

### Installation

1. **Clone or download this project:**
   ```bash
   git clone <your-repo-url>
   cd hoctap-api-project
   ```

2. **Configure Environment Variables:**
   
   Copy the example configuration:
   ```bash
   cp env.example config.env
   ```
   
   Edit `config.env` with your database credentials:
   ```env
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=hoctap_user
   DB_PASSWORD=your_secure_password
   DB_NAME=hoctap_api

   # Server Configuration
   SERVER_PORT=8080

   # Environment
   ENVIRONMENT=development
   ```

3. **Install dependencies:**
   ```bash
   go mod tidy
   ```

4. **Run the server:**
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080` and automatically:
- Connect to MySQL database
- Create necessary tables
- Seed initial test data

## API Endpoints

### Base Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | HTML dashboard |
| GET | `/health` | Health check with database status |
| GET | `/welcome` | API welcome message |

### User Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users` | Get all users |
| GET | `/api/users/{id}` | Get user by ID |
| POST | `/api/users` | Create a new user |
| PUT | `/api/users/{id}` | Update user by ID |
| DELETE | `/api/users/{id}` | Delete user by ID |
| GET | `/api/users/stats` | Get user statistics |

### Example Requests

#### Get all users
```bash
curl http://localhost:8080/api/users
```

#### Get user by ID
```bash
curl http://localhost:8080/api/users/1
```

#### Create a new user
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Johnson", "email": "alice@example.com"}'
```

#### Update a user
```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "John Smith", "email": "johnsmith@example.com"}'
```

#### Delete a user
```bash
curl -X DELETE http://localhost:8080/api/users/1
```

#### Get user statistics
```bash
curl http://localhost:8080/api/users/stats
```

#### Health check
```bash
curl http://localhost:8080/health
```

## Response Format

All API responses follow this standard format:

```json
{
  "message": "Success message",
  "data": {}, 
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Development

### Project Structure

```
hoctap-api-project/
├── main.go              # Main application file
├── database/            # Database layer
│   ├── connection.go    # Database connection management
│   └── user.go         # User model and repository
├── config.env          # Environment configuration
├── env.example         # Example environment file
├── index.html          # HTML dashboard
├── styles.css          # Dashboard styling
├── script.js           # Dashboard JavaScript
├── go.mod              # Go module dependencies
├── go.sum              # Dependency checksums
├── start-servers.bat   # Server launcher script
└── README.md           # This file
```

### Environment Configuration

The application uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | MySQL host | `localhost` |
| `DB_PORT` | MySQL port | `3306` |
| `DB_USER` | MySQL username | `root` |
| `DB_PASSWORD` | MySQL password | `` |
| `DB_NAME` | Database name | `hoctap_api` |
| `SERVER_PORT` | Server port | `8080` |
| `ENVIRONMENT` | Environment mode | `development` |

### Running in Development

For development with auto-reload, you can use:

```bash
# Install air for live reloading (optional)
go install github.com/cosmtrek/air@latest

# Run with air
air
```

Or simply run:
```bash
go run main.go
```

### Building for Production

```bash
# Build binary
go build -o hoctap-api main.go

# Run binary
./hoctap-api
```

### Database Schema

The application automatically creates the following table:

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
```

### Running with Docker (Optional)

If you prefer to use Docker for MySQL:

```bash
# Start MySQL container
docker run --name mysql-hoctap \
  -e MYSQL_ROOT_PASSWORD=rootpassword \
  -e MYSQL_DATABASE=hoctap_api \
  -e MYSQL_USER=hoctap_user \
  -e MYSQL_PASSWORD=your_secure_password \
  -p 3306:3306 \
  -d mysql:8.0

# Update config.env accordingly
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Verify MySQL is running: `mysql -u root -p`
   - Check credentials in `config.env`
   - Ensure database exists

2. **Port Already in Use**
   - Change `SERVER_PORT` in `config.env`
   - Or stop the process using the port

3. **Permission Denied**
   - Ensure database user has proper privileges
   - Check MySQL user permissions

### Logs

The application logs important events:
- Database connection status
- API requests with timing
- Error messages with details

## License

This project is open source and available under the [MIT License](LICENSE).