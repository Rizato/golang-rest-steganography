FROM golang:trixie AS builder
WORKDIR /build
COPY . .
RUN mkdir bin

RUN go build -o bin/worker ./cmd/worker
RUN go build -o bin/server ./cmd/server
RUN go build -o bin/embed ./cmd/embed
RUN go build -o bin/extract ./cmd/extract


FROM debian:trixie

RUN apt-get -y update && apt-get install -y \
    libc6 \
    libstdc++6 \
    libprotobuf32 \
    libnl-route-3-200

WORKDIR /app
COPY --from=builder /build/bin/worker .
COPY --from=builder /build/bin/server .
COPY --from=builder /build/bin/embed .
COPY --from=builder /build/bin/extract .
COPY --from=nsjailcontainer /bin/nsjail .

COPY steg.cfg .

CMD ["/app/server"]
