FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /3rike-api .

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /3rike-api .
EXPOSE 8080
CMD ["./3rike-api"]
