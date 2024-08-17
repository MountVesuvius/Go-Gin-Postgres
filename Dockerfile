# Stage 1: Build the Go binary - Go 1.22
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copy Go module files, download dependencies, copy application code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the Go binary
RUN GOOS=linux GOARCH=amd64 go build -o main .
# Might need to statically link things if prod is running on scratch
# CGO_ENABLED=0

# Stage 2: Development Environment with Hot Reload - Go 1.22
FROM golang:1.22-alpine AS dev
WORKDIR /app

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Copy Go module files, download dependencies, copy application code
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Run the app with air for live reloading
CMD ["air"]

# Stage 3: Production Environment with Scratch (minimal)
FROM scratch AS prod
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080
EXPOSE 8080

# Command to run the compiled binary
CMD ["./main"]
