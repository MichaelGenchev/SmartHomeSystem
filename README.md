# Smart Home Management System

## Overview

The Smart Home Management System is a microservices-based application designed to control and monitor various smart devices in a home environment. This system provides a scalable and flexible architecture for managing smart home devices, user accounts, and device states.

## Table of Contents

1. [Architecture](#architecture)
2. [Services](#services)
3. [Technologies Used](#technologies-used)
4. [Project Structure](#project-structure)
5. [Getting Started](#getting-started)
6. [Development](#development)
7. [Testing](#testing)
8. [Deployment](#deployment)
9. [Contributing](#contributing)
10. [License](#license)

## Architecture

The Smart Home Management System follows a microservices architecture, with the following key components:

```
                    ┌─────────────────┐
                    │   API Gateway   │
                    └─────────────────┘
                             │
           ┌─────────────────┼─────────────────┐
           │                 │                 │
  ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
  │  Device Service  │ │   User Service   │ │  Other Services  │
  └──────────────────┘ └──────────────────┘ └──────────────────┘
           │                 │                 │
  ┌──────────────────┐ ┌──────────────────┐ ┌──────────────────┐
  │     MongoDB      | │  PostgreSQL      │ │  Other Databases │
  └──────────────────┘ └──────────────────┘ └──────────────────┘
```

- **API Gateway**: Routes requests to appropriate microservices.
- **Device Service**: Manages smart home devices, their states, and operations.
- **User Service**: Handles user authentication, authorization, and profile management.
- **Other Services**: Placeholder for future services (e.g., Automation Service, Notification Service).
- **Databases**: Each service can have it's own database.

## Services

### Device Service

The Device Service is responsible for managing smart home devices. It provides the following functionalities:

- Create, read, update, and delete devices
- Manage device states
- List devices for a user
- Implement CQRS pattern with separate read and write models

### User Service

The User Service handles user-related operations, including:

- User registration and authentication
- User profile management
- Authorization and access control

## Technologies Used

- **Go**: Primary programming language
- **gRPC**: For inter-service communication
- **Protocol Buffers**: For defining service contracts
- **Docker**: For containerization
- **Kubernetes**: For orchestration and deployment
- **PostgreSQL**: Relational database for User Service
- **MongoDB**: NoSQL database for Device Service
- **Makefile**: For automating build, test, and deployment tasks
- **GitHub Actions**: For CI/CD pipelines

## Project Structure

```
smart-home-system/
├── cmd/
│   ├── device-service/
│   └── user-service/
├── internal/
│   ├── device/
│   └── user/
├── pkg/
│   ├── models/
│   └── proto/
├── deployments/
│   ├── docker-compose.yml
│   └── kubernetes/
├── build/
│   ├── Dockerfile.device
│   └── Dockerfile.user
├── .github/workflows/
├── Makefile
├── go.mod
└── README.md
```

## Getting Started

To get started with the Smart Home Management System:

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/smart-home-system.git
   ```

2. Install dependencies:
   ```
   make deps
   ```

3. Build the services:
   ```
   make build
   ```

4. Run the services locally:
   ```
   make run-device
   make run-user
   ```

## Development

To develop new features or fix bugs:

1. Create a new branch:
   ```
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and write tests.

3. Run tests:
   ```
   make test
   ```

4. Run integration tests:
   ```
   make integration-test-all
   ```

5. Submit a pull request for review.

## Testing

- Run unit tests: `make test`
- Run integration tests: `make integration-test-all`
- Generate coverage report: `make coverage`

## Deployment

To deploy the Smart Home Management System:

1. Build Docker images:
   ```
   make docker-build-all
   ```

2. Push Docker images:
   ```
   make docker-push-all
   ```

3. Deploy to Kubernetes:
   ```
   make k8s-deploy-all
   ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
