# Use the official Golang image as the base
FROM golang:1.22.4

# Time Zone
ENV TZ=Asia/Jakarta

# Set the working directory inside the container
WORKDIR /app

# Copy the Go application source code to the container
COPY . .

# Build the Go application
RUN go mod tidy
RUN go build -o sigmatech-test

# Set the environment variable for the port
ENV PORT 9000

# Expose the port the application will run on
EXPOSE 9000

# Command to run the Go application
CMD ["./sigmatech-test"]
