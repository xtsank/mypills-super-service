FROM golang:1.26.2-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /mypills-app ./src/cmd/app/main.go ./src/cmd/app/app.go

FROM alpine:latest
WORKDIR /root/

COPY --from=builder /mypills-app .

COPY .env .

COPY docs/ ./docs/

CMD ["./mypills-app"]