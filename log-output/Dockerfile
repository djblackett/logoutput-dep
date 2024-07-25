FROM golang:latest as builder
LABEL authors="dave"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /random-letters

FROM alpine:latest

COPY --from=builder /random-letters .
CMD ["/random-letters"]