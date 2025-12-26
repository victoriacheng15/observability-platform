FROM golang:1.25-alpine

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the root .env file to the current directory (which is /app)
COPY .env .

# Copy go.mod and go.sum from the proxy directory
COPY proxy/go.mod proxy/go.sum ./
RUN go mod download

# Copy the rest of the proxy source code
COPY proxy/ .

RUN go build -o proxy_server .

EXPOSE 8085

CMD ["./proxy_server"]