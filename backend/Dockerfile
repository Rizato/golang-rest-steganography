FROM golang:alpine AS builder
WORKDIR /build
COPY . .

RUN go build

FROM alpine
WORKDIR /app
COPY --from=builder /build/steg .

CMD ["/app/steg"]