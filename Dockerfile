FROM golang:1.21-alpine as build

WORKDIR /src
COPY . .

RUN go build -o /build/author cmd/author/main.go 

FROM alpine:latest
WORKDIR /app
COPY --from=build /build/author .

RUN adduser -D otel
USER otel:otel

