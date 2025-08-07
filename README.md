# Steganography API Service

A Go-based REST API for hiding messages within images using steganography techniques. This service allows you to encode secret messages into JPG and PNG images, and decode them back.

## Overview

Steganography is the practice of concealing messages or information within other non-secret data. Unlike cryptography which makes messages unreadable, steganography makes messages undetectable. This service embeds text messages into image files that appear completely normal but contain hidden data.

## Features

- **Image Message Encoding**: Hide text messages within JPG and PNG images
- **Image Message Decoding**: Extract hidden messages from steganographic images  
- **Async Processing**: Non-blocking job processing with status tracking
- **RESTful API**: Clean HTTP endpoints for all operations
- **Middleware Stack**: Request logging, CORS, and statistics tracking
- **Generic Handler Pattern**: Type-safe, reusable handlers using Go generics

## Project Structure

```
.
├── main.go                 # Application entry point
├── go.mod                  # Go module dependencies
├── go.sum                  # Dependency checksums
├── LICENSE                 # Project license
├── README.md               # This file
├── api/
│   └── v1/
│       ├── handlers/       # HTTP request handlers
│       │   ├── factory.go  # Job factory implementations
│       │   └── handlers.go # Generic and specific handlers
│       ├── models/         # Data models and interfaces
│       │   └── models.go   # Job models and status enums
│       └── views/          # Route configuration
│           └── views.go    # API endpoint routing
├── files/                  # File handling utilities
│   └── files.go            # File validation and processing
└── middleware/             # HTTP middleware
    ├── middleware.go       # Middleware configuration
    ├── log.go              # Request logging
    └── stats.go            # Statistics collection
```

## Installation

### Prerequisites

- Go 1.24 or higher
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

#### Encode Message into Image

```http
POST /api/v1/encode
```
**Request:**
- Method: `POST`
- Content-Type: `multipart/form-data`
- Form Fields:
  - `image`: Image file (JPG or PNG)
  - `message`: Text message to hide

**Response:**
```json
{
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "status": "Submitted",
  "status-message": "",
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:00Z"
}
```

#### Get Encode Job Status

```http
GET /api/v1/encode/{uuid}
```
**Response:**
```json
{
  "uuid": "550e8400-e29b-41d4-a716-446655440000",
  "status": "Complete",
  "status-message": "Image processed successfully",
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:01Z"
}
```

#### Decode Message from Image

```http
POST /api/v1/decode
```
**Request:**
- Method: `POST`
- Content-Type: `multipart/form-data`
- Form Fields:
  - `image`: Steganographic image file

**Response:**
```json
{
  "uuid": "660e8400-e29b-41d4-a716-446655440000",
  "status": "Submitted",
  "status-message": "",
  "message": "",
  "createdAt": "2024-01-01T12:00:00Z",
  "updatedAt": "2024-01-01T12:00:00Z"
}
```

#### Get Decode Job Status

```http
GET /api/v1/decode/{uuid}
```
**Response:**
```json
{
  "uuid": "660e8400-e29b-41d4-a716-446655440000",
  "status": "Complete",
  "status-message": "Message extracted successfully",
  "message": "This is the hidden message",
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
    "GET /api/v1/encode/{uuid}": {
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
- `Error` (0): Processing failed
- `Submitted` (1): Job received and queued
- `InProgress` (2): Currently processing
- `Complete` (3): Successfully completed

## Data Models

### EncodeJob
```go
type EncodeJob struct {
    Uuid          uuid.UUID
    Status        Status
    StatusMessage string
    ImagePath     string    // Internal use only
    StegImagePath string    // Internal use only
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
```

### DecodeJob
```go
type DecodeJob struct {
    Uuid          uuid.UUID
    Status        Status
    StatusMessage string
    ImagePath     string    // Internal use only
    Message       string
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
- [x] Generic job handlers with type safety
- [x] Request logging middleware
- [x] Statistics middleware
- [x] File upload handling
- [x] In-memory Job status tracking

### Planned Features
- [ ] Actual steganography implementation
- [ ] Asynchronous image processing with goroutines
- [ ] S3 integration for image storage
- [ ] PostgreSQL for persistent job storage
- [ ] Message queue with RabbitMQ
- [ ] Protobuf serialization
- [ ] Security sandboxing with nsjail
- [ ] Web frontend
- [ ] Authentication and authorization
- [ ] CORS configuration
- [ ] Automatic file cleanup
- [ ] Kubernetes deployment
- [ ] AWS infrastructure

## License

This project is licensed under the terms in the LICENSE file.

