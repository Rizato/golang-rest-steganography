# Steganography API Service

A Go-based REST API for hiding messages within images using steganography techniques. This service allows you to embed
secret messages into JPG and PNG images, and extract them later.

## ⚠️ Important Notice

**This is a learning project** created to explore Go and backend development concepts.
It is NOT production-ready and should not be deployed publicly without understanding the significant security
implications of running an image hosting service.
Image hosting services can be exploited for storing and distributing malicious content, illegal materials, and more.

## Overview

Steganography is the practice of concealing messages or information within other non-secret data. Unlike cryptography
which makes messages unreadable, steganography makes messages undetectable. This service embeds text messages into image
files that appear completely normal but contain hidden data.

## Features

- **Image Message Embedding**: Embed text messages within JPG and PNG images
- **Image Message Extracting**: Extract hidden messages from steganographic images
- **Async Processing**: Non-blocking job processing with RabbitMQ message queue
- **Security Sandboxing**: NsJail isolation for steganography operations
- **RESTful API**: Clean HTTP endpoints with Go 1.24 pattern matching
- **Web Interface**: Vue.js frontend for easy interaction
- **Database Persistence**: PostgreSQL for reliable job and file storage
- **Middleware Stack**: Request logging, CORS, and statistics tracking
- **Generic Handler Pattern**: Type-safe, reusable handlers using Go generics
- **Docker Deployment**: Containerized architecture for easy deployment

## Project Structure

```
.
├── docker-compose.yml       # Docker services configuration
├── LICENSE                  # Project license
├── README.md                # This file
├── backend/                 # Backend Go service
│   ├── Dockerfile           # Backend container configuration
│   ├── go.mod               # Go module dependencies
│   ├── go.sum               # Dependency checksums
│   ├── steg.cfg             # NsJail security configuration
│   ├── cmd/                 # Application entry points
│   │   ├── server/          # REST API server
│   │   └── worker/          # Async job worker
│   ├── server/              # HTTP server implementation
│   │   ├── handlers/        # HTTP request handlers
│   │   ├── middleware/      # HTTP middleware
│   │   └── views/           # Route configuration
│   ├── worker/              # Job processing worker
│   │   └── processor/       # Job processing logic
│   ├── shared/              # Shared components
│   │   ├── models/          # Data models
│   │   ├── datastore/       # Database interactions
│   │   └── rabbitmq/        # Message queue client
│   ├── steganography/       # Steganography implementation
│   │   └── stegtool/        # Core steganography logic
│   └── migrations/          # Database migrations
└── frontend/                # Vue.js web interface
    ├── src/                 # Frontend source code
    ├── public/              # Static assets
    └── package.json         # Node.js dependencies
```

## Installation

### Prerequisites

- Go 1.24 or higher (required for URL pattern matching)
- Docker and Docker Compose
- Git

### Setup

1. Clone the repository:

```bash
git clone git@github.com:Rizato/golang-rest-steganography.git
cd golang-rest-steganography
```

2. Start all services with Docker Compose:

```bash
docker-compose up --build
```

This will start:
- PostgreSQL database on port 5432
- RabbitMQ message queue on ports 5672 (AMQP) and 15672 (Management UI)
- Backend API server on port 8080
- Frontend web interface on port 3000
- Background worker for processing jobs

3. Access the application:
   - Web Interface: http://localhost:3000
   - API: http://localhost:8080/api/v1
   - RabbitMQ Management: http://localhost:15672 (guest/guest)

## API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Endpoints

#### Upload Image

```http
POST /api/v1/images
```

**Request:**

- Method: `POST`
- Content-Type: `multipart/form-data`
- Form Fields:
    - `file`: Image file (JPG or PNG)

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "filename": "image.jpg",
  "size": 102400,
  "createdAt": "2024-01-01T12:00:00Z"
}
```

#### Download Image

```http
GET /api/v1/images/{id}/download
```

**Response:**

- Binary image data with appropriate Content-Type header

#### Get/Delete Image

```http
GET /api/v1/images/{id}
DELETE /api/v1/images/{id}
```

**GET Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "filename": "image.jpg",
  "size": 102400,
  "createdAt": "2024-01-01T12:00:00Z"
}
```

#### Create Embed Job

```http
POST /api/v1/embed
```

**Request:**

```json
{
  "imageId": "550e8400-e29b-41d4-a716-446655440000",
  "message": "Secret message to hide"
}
```

**Response:**

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "status": "Submitted",
  "statusMessage": "",
  "imageId": "550e8400-e29b-41d4-a716-446655440000",
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:00Z"
}
```

#### List Embed Jobs

```http
GET /api/v1/embed
```

**Response:**

```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "status": "Complete",
    "statusMessage": "Processing complete",
    "imageId": "550e8400-e29b-41d4-a716-446655440000",
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:00:01Z"
  }
]
```

#### Get/Delete Embed Job

```http
GET /api/v1/embed/{id}
DELETE /api/v1/embed/{id}
```

**GET Response:**

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "status": "Complete",
  "statusMessage": "Processing complete",
  "imageId": "550e8400-e29b-41d4-a716-446655440000",
  "outputImageId": "770e8400-e29b-41d4-a716-446655440000",
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:01Z"
}
```

#### Start Embed Job Processing

```http
POST /api/v1/embed/{id}/start
```

**Response:**

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "status": "In Progress",
  "statusMessage": "Processing started",
  "imageId": "550e8400-e29b-41d4-a716-446655440000",
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:01Z"
}
```

#### Create Extract Job

```http
POST /api/v1/extract
```

**Request:**

```json
{
  "imageId": "770e8400-e29b-41d4-a716-446655440000"
}
```

**Response:**

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "status": "Submitted",
  "statusMessage": "",
  "imageId": "770e8400-e29b-41d4-a716-446655440000",
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:00Z"
}
```

#### List Extract Jobs

```http
GET /api/v1/extract
```

**Response:**

```json
[
  {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "status": "Complete",
    "statusMessage": "Extraction complete",
    "imageId": "770e8400-e29b-41d4-a716-446655440000",
    "message": "Secret message to hide",
    "createdAt": "2024-01-01T12:00:00Z",
    "updatedAt": "2024-01-01T12:00:01Z"
  }
]
```

#### Get/Delete Extract Job

```http
GET /api/v1/extract/{id}
DELETE /api/v1/extract/{id}
```

**GET Response:**

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "status": "Complete",
  "statusMessage": "Extraction complete",
  "imageId": "770e8400-e29b-41d4-a716-446655440000",
  "message": "Secret message to hide",
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:01Z"
}
```

#### Start Extract Job Processing

```http
POST /api/v1/extract/{id}/start
```

**Response:**

```json
{
  "id": "880e8400-e29b-41d4-a716-446655440000",
  "status": "In Progress",
  "statusMessage": "Processing started",
  "imageId": "770e8400-e29b-41d4-a716-446655440000",
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:01Z"
}
```

#### Statistics

```http
GET /api/v1/stats
```

**Response:**

```json
{
  "total_requests": 42,
  "uptime": "4.011104257s",
  "endpoints": {
    "GET /api/v1/embed/{uuid}": {
      "count": 10,
      "total_time": "963.91µs",
      "avg_time": "96.391µs"
    },
    "GET /api/v1/stats": {
      "count": 5,
      "total_time": "376.88µs",
      "avg_time": "75.376µs"
    }
  }
}
```

### Status Values

Jobs can have the following status values:

- `Error`: Processing failed
- `Submitted`: Job created and ready to process
- `In Progress`: Currently processing
- `Complete`: Successfully completed
- `Cancelled`: Job was cancelled

## Data Models

### ServerFile

```go
type ServerFile struct {
ID        uuid.UUID
Filename  string
Size      int64
Path      string // Internal use only
CreatedAt time.Time
}
```

### EmbedJob

```go
type EmbedJob struct {
ID            uuid.UUID
Status        Status
StatusMessage string
ImageID       uuid.UUID
OutputImageID uuid.UUID // Set after processing
Message       string
CreatedAt     time.Time
UpdatedAt     time.Time
}
```

### ExtractJob

```go
type ExtractJob struct {
ID            uuid.UUID
Status        Status
StatusMessage string
ImageID       uuid.UUID
Message       string // Extracted message
CreatedAt     time.Time
UpdatedAt     time.Time
}
```

## Development

### Running Tests

```bash
cd backend
go test ./...
```

### Building with Docker

```bash
docker-compose build
```

### Running Individual Services

For development, you can run services individually:

```bash
# Start infrastructure services only
docker-compose up postgres rabbitmq

# Run backend locally
cd backend
go run cmd/server/main.go

# Run worker locally
cd backend
go run cmd/worker/main.go

# Run frontend locally
cd frontend
npm install
npm run dev
```

## Roadmap

### Current Implementation

- [x] Basic REST API structure
- [x] Generic CRUD handlers with type safety
- [x] Request logging middleware
- [x] Statistics middleware
- [x] File upload and download handling
- [x] Thread-safe in-memory datastore
- [x] Separate job creation and processing endpoints
- [x] Request validation
- [x] PostgreSQL for persistent job storage
- [x] Web frontend (Created with claude)
- [x] Actual steganography implementation (created with claude)
- [x] Asynchronous image processing with goroutines
- [x] CORS configuration
- [x] Message queue with RabbitMQ
- [x] Security sandboxing with nsjail

### Planned Features

- [ ] S3 integration for image storage
- [ ] Protobuf serialization
- [ ] Authentication and authorization
- [ ] Automatic file cleanup
- [ ] Kubernetes deployment
- [ ] AWS infrastructure

## License

This project is licensed under the terms in the LICENSE file.
