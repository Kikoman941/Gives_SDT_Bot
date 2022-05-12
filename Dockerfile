FROM golang:1.18-alpine as builder
WORKDIR /build
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o /app cmd/main.go

FROM alpine:3.15.4
COPY .env /.env
COPY --from=builder app /app
ENTRYPOINT ["/app"]