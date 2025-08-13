## Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN apk add --no-cache git
COPY . .
RUN go mod download
RUN go mod tidy
RUN go build -o bank-api .

## Run stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/bank-api ./bank-api
EXPOSE 8080
ENTRYPOINT ["/app/bank-api"]


