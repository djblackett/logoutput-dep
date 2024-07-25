FROM golang:latest as builder
LABEL authors="djblackett"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /ping-pong

FROM scratch

COPY --from=builder /ping-pong .
CMD ["/ping-pong"]