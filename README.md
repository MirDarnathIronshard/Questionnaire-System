
# Voting System

A highly scalable and modular voting system designed to support various types of questionnaires, voting mechanisms, and user management. Built with **Go**, this project adheres to modern software engineering practices.

## Features
- User authentication and authorization
- Comprehensive questionnaire and voting management
- Real-time notifications and updates
- Modular architecture for scalability and maintainability
- Extensive documentation and clean codebase
- Integration and unit testing for core functionality

## Project Structure
```plaintext
Voting-System/
│
├── docker/                # Docker configurations and setup files
├── .git/                  # Git version control files
├── .idea/                 # IDE-specific files (e.g., IntelliJ/PyCharm)
├── src/                   # Main source code directory
│   ├── cmd/               # Entry point of the application
│   ├── config/            # Configuration files
│   ├── internal/          # Core business logic
│   │   ├── middleware/    # Middlewares for logging, security, etc.
│   │   ├── repositories/  # Database access and queries
│   │   ├── models/        # Data models
│   │   ├── services/      # Service layer
│   │   ├── handlers/      # API request handlers
│   │   └── routes/        # API routes definition
│   ├── pkg/               # Reusable utility packages
│   ├── tests/             # Unit and integration tests
│   └── docs/              # Documentation files
├── go.mod                 # Go modules file for dependency management
├── go.sum                 # Checksums for module dependencies
└── .env.example           # Sample environment configuration file
```

## Services Overview
The system is divided into various services, each handling a specific domain of the application:

### 1. **Authentication Service**
- Handles user login, registration, and token-based authentication.
- Middleware integration for request authorization.
- Implements password hashing and token validation.

### 2. **Questionnaire Service**
- Manages the creation, update, and deletion of questionnaires.
- Supports user-defined roles and permissions for accessing questionnaires.

### 3. **Voting Service**
- Handles the submission, storage, and validation of votes.
- Provides real-time updates for active voting sessions.

### 4. **Notification Service**
- Sends real-time notifications to users about updates, reminders, or errors.
- Uses websockets for live notifications and email for asynchronous messaging.

### 5. **User Management Service**
- Manages user profiles, roles, and access permissions.
- Tracks user activity and access to different modules.

### 6. **Chat Service**
- Enables real-time messaging for collaborative voting or discussions.
- Stores chat history and provides searchable interfaces.

### 7. **Questionnaire Role Management**
- Assigns and verifies roles for users participating in specific questionnaires.
- Supports fine-grained access control for questionnaire operations.

### 8. **Error Handling Service**
- Centralized error handling for consistent API responses.
- Maintains an error code library to standardize error management.

### 9. **Logger Service**
- Provides structured logging for debugging and monitoring.
- Supports different logging levels (info, warning, error).

## ERD (Entity-Relationship Diagram)
The database design is visualized in the ERD, which illustrates the relationships between the main entities of the system, including:
- Users
- Questionnaires
- Questions
- Responses
- Votes
- Roles and Permissions

The ERD is available in the following formats:
- **Markdown file**: `src/docs/erd.puml`
- **PNG image**: `src/docs/erd.png`

### Preview of the ERD:
The diagram demonstrates how entities interact and their relational structure, ensuring a robust database design for scalability and performance.

> For detailed analysis, open the PNG file in the `docs` folder.

## Getting Started
Follow these steps to set up and run the project locally.

### Prerequisites
- Go (version 1.20+ recommended)
- Docker (optional, for containerized environments)
- PostgreSQL (or any supported database)

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/voting-system.git
   cd voting-system
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Set up environment variables:
    - Copy `.env.example` to `.env` and update values as needed.

4. (Optional) Run with Docker:
   ```bash
   docker-compose up --build
   ```

5. Start the application:
   ```bash
   go run src/cmd/main.go
   ```

### Running Tests
Execute the following command to run all tests:
```bash
go test ./...
```

## API Documentation
Refer to the API documentation located in `src/docs/api_documentation.md`.

## Contributing
We welcome contributions! Please follow these steps:
1. Fork the repository.
2. Create a new branch for your feature/bugfix.
3. Submit a pull request with a clear description of your changes.

## License
This project is licensed under the MIT License. See the `LICENSE` file for details.

 