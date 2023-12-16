FROM golang:1.21-alpine AS build

# Set the working directory inside the container
WORKDIR /Kliptopia

# Copy the Go application source code to the container
COPY . .

# Build the Go application
RUN cd cmd && go build -o cmd/main .

FROM alpine:latest

# Set the working directory inside the container
WORKDIR /Kliptopia

# Copy only the built binary
COPY --from=build /Kliptopia/cmd/main .

# Expose the port 
EXPOSE 9000

# Command to run the executable
CMD ["./main"]
