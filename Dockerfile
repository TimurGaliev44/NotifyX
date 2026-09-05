FROM golang:1.27-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /check cmd/check/main.go

FROM alpine:3.19
COPY --from=builder /check .

CMD ["./check"]
