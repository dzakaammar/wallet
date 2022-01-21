FROM golang:alpine AS builder

RUN apk update && apk upgrade && apk add --no-cache git tzdata ca-certificates

WORKDIR /app
COPY . .

RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -ldflags="-s -w" -o binary


FROM scratch

WORKDIR /app
COPY --from=builder /app/binary /app/binary
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8000

ENTRYPOINT ["/app/binary"]