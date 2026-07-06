FROM golang:1.26.4-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/app/exe /app/cmd/app

FROM alpine:3.23
RUN apk add --no-cache tzdata
WORKDIR /app
COPY --from=builder /app/cmd/app/exe /app
ENV PROJECT_ROOT=/app
CMD ["/app/exe"]