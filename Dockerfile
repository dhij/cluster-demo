FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./userapi

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /build/userapi ./userapi
CMD ["/app/userapi"]