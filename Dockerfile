# Use the specific Golang image as the base image
FROM golang:1.22-rc-alpine3.18 as builder

WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o bubblr_server cmd/main.go

#FROM scratch
FROM centos:latest
COPY --from=builder /app/bubblr_server /

# Copy CA certificates to prevent x509: certificate signed by unknown authority errors
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 8081
ENTRYPOINT ["./bubblr_server"]