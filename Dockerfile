FROM golang:1.21-alpine as build

WORKDIR /src
COPY . .

RUN go build -o /build/author cmd/author/main.go 
RUN go build -o /build/book cmd/book/main.go 

FROM alpine:latest
WORKDIR /app
COPY --from=build /build/author .
COPY --from=build /build/book .

RUN adduser -D otel
USER otel:otel

