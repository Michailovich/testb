FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /testb

FROM alpine:latest

WORKDIR /app

COPY --from=builder /testb /app/testb
COPY gql/schema.graphql /app/gql/schema.graphql

EXPOSE 8080

CMD ["/app/testb"]