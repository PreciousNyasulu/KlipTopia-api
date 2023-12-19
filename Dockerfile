FROM golang:1.21-alpine AS build

# Set the working directory inside the container
WORKDIR /Kliptopia

# Copy the Go application source code to the container
COPY . .

# Build the Go application
RUN go build -o build/main cmd/main.go

FROM alpine:latest

# Set the working directory inside the container
WORKDIR /Kliptopia

# Copy only the built binary
COPY --from=build /Kliptopia/build/main .

# Expose the port 
EXPOSE 9000

# Command to run the executable
CMD ["./main"]
