FROM golang:alpine AS builder
WORKDIR /build
COPY . .
RUN mkdir bin

RUN go build -o bin/worker ./cmd/worker
RUN go build -o bin/server ./cmd/server
RUN go build -o bin/stegtool ./cmd/stegtool


FROM alpine
WORKDIR /app
COPY --from=builder /build/bin/worker .
COPY --from=builder /build/bin/server .
COPY --from=builder /build/bin/stegtool .

CMD ["/app/server"]
