# Use official Golang image as the base
FROM golang:1.23.4 AS build
WORKDIR /app
# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy the source code
COPY . .
# Build the application with static linking
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gin-api .

# Use a minimal base image to run the app
FROM alpine:latest
WORKDIR /app
# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates
# Copy the compiled binary from the build stage
COPY --from=build /app/gin-api .
# Expose API port
EXPOSE 8080
# Run the application
CMD ["./gin-api"]