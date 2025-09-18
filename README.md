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
- **Async Processing**: Non-blocking job processing with status tracking
- **RESTful API**: Clean HTTP endpoints for all operations
- **Middleware Stack**: Request logging, CORS, and statistics tracking
- **Generic Handler Pattern**: Type-safe, reusable handlers using Go generics

## Project Structure

```
.
├── main.go                  # Application entry point
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency checksums
├── LICENSE                  # Project license
├── README.md                # This file
├── api/ 
│   └── v1/ 
│       ├── handlers/        # HTTP request handlers
│       │   ├── crud.go      # CRUD operations
│       │   ├── handlers.go  # Job-specific handlers
│       │   ├── requests.go  # Request processing logic
│       │   └── uploads.go   # File upload handling
│       ├── models/          # Data models and storage
│       │   ├── datastore.go # Thread-safe job storage
│       │   └── models.go    # Job models and status enums
│       └── views/           # Route configuration
│           └── views.go     # API endpoint routing
├── crud/                    # CRUD interfaces
│   └── crud.go              # Generic CRUD operations
├── files/                   # File handling utilities
│   └── files.go             # File validation and processing
└── middleware/              # HTTP middleware
    ├── middleware.go        # Middleware configuration
    ├── log.go               # Request logging
    └── stats.go             # Statistics collection
```

## Installation

### Prerequisites

- Go 1.24 or higher
-
- Git

### Setup

1. Clone the repository:

```bash
git clone git@github.com:Rizato/golang-rest-steganography.git
cd golang-rest-steganography
```

2. Install dependencies:

```bash
go mod download
```

3. Run the application:

```bash
go run main.go
```

The server will start on port 8080.

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
go test ./...
```

### Building

```bash
go build -o steg
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
