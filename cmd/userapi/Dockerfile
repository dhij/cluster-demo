FROM golang:1.21-alpine AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./userapi ./cmd/userapi/main.go

RUN CGO_ENABLED=0 GOOS=linux go install -a -ldflags "-s -w" -tags "mysql" \
    github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0

FROM gcr.io/distroless/base-debian12

WORKDIR /app
ENV PATH=/app/bin:$PATH

COPY --from=builder /build/userapi ./bin/userapi
COPY --from=builder /go/bin/migrate ./bin/migrate
COPY --from=builder /build/db/ .

ENTRYPOINT ["userapi"]