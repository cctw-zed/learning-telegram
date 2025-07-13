# Stage 1: Build the Go binary
# Use our custom builder image which has build-base and goproxy pre-configured.
FROM telegram-builder:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application with CGO explicitly enabled.
# The build tools are pre-installed in the 'telegram-builder' image.
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o /server cmd/server/main.go


# Stage 2: Create the final, minimal production image
FROM alpine:latest

# Install SQLite runtime libraries, which are required by the cgo-enabled binary
RUN apk add --no-cache sqlite-libs

# Copy the built binary to a standard system path
COPY --from=builder /server /usr/local/bin/server

# Create a dedicated directory for persistent data
RUN mkdir /data

# Set the working directory to the data directory
WORKDIR /data

# Expose the port the application runs on
EXPOSE 8080

# The command to run the application. 'server' is found in the PATH.
# The database file 'telegram.db' will be created in the current working
# directory ('/data'), which is the volume mount point.
CMD ["server"] 