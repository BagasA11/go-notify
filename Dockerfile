FROM golang:1.23-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the Go application
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
