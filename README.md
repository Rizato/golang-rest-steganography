# Go Learning

I built this project to learn Go. 

I figured I would build a fairly simple backend using net/http, goroutines, and a database. 

## Structure

```
main.go
| handlers
    handlers.go
| middleware
    middleware.go 
```

## Steganography

Steganography is a concept in hiding messages. 
Where cryptography makes the message hard to read, steganography makes the message hard to find.

In this case, we are embedding text into jpg, and png images. 
You will be able to host the image, or send it via text, and it will have a hidden message.

Your recipient will have to know to go to this service to extract the message

## Goals

- Build a simple RESTful backend with asynchronous image processing
- Use goroutine to upload the image to S3 (or save locally for first pass)
- Use goroutines to process the image (insecure, but to learn goroutines, and channels)
- Use PostgreSQL to store job state
- Replace goroutines with microservice using rabbitmq and protobuf
- Use nsjail on encode/decode service to protect from malicious user uploads
- Use claude to make a simple frontend
- Middleware (csrf, logging, auth?)
- Delete images after use
- Deploy on k8s on aws


## Models

- Encode
  - uuid
  - status
  - rawImageUrl
  - stegImageUrl
  - uploadConsent
  - uploadConsentedAt
- Decode
  - uuid
  - status
  - stegImageUrl
  - message

## Views

- POST /api/v1/encode
- GET /api/v1/encode/{id}
- POST /api/v1/decode
- GET /api/v1/decode/{id}

