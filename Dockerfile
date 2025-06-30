# -----------------------
# Stage 1: Build the Go app
# -----------------------
FROM golang:1.24 AS builder

WORKDIR /Job-Service

# Copy go module files and download deps early for better caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy rest of the code
COPY . .

# Build the Go binary
RUN go build -o job-service ./cmd/job-service

# -----------------------
# Stage 2: Runtime (minimal image)
# -----------------------
FROM debian:bookworm-slim  

WORKDIR /Job-Service

# Optional: Create logs dir
RUN mkdir -p logs

# Copy the binary from the builder stage
COPY --from=builder /Job-Service/job-service .

# Expose the port your app listens on
EXPOSE 8080

# Command to run the binary
CMD ["./job-service"]
