# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy the go directory from your project root into the container
COPY go/ .

# Initialize the module if go.mod is missing, otherwise tidy up
# This handles the case where go.mod might not be in the screenshot but exists
RUN [ ! -f go.mod ] && go mod init au_portal_bot || true
RUN go mod tidy

# Build the binary with size optimizations (strip debug info)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o au_bot main.go

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /build/au_bot .

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Run the binary
CMD ["./au_bot"]
