FROM golang:1.21-alpine AS builder

WORKDIR /code

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o echo-server .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /code/echo-server .

CMD ["./echo-server"]