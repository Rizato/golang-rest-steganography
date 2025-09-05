FROM golang:alpine AS build
WORKDIR /build
COPY .. .

RUN go build

FROM alpine
WORKDIR /app
COPY --from=build /build/steg .

CMD ["/app/steg"]