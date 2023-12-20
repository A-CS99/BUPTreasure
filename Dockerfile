# Use the official Golang image as the base image
FROM golang:1.21.5

# Set the working directory
WORKDIR /app

# Copy the local code to the container
COPY . .

# 设置 Go 代理
ENV GOPROXY=https://goproxy.cn,direct

# 设置 GIN_MODE 为 "release"
ENV GIN_MODE=release

# Build the Go application
RUN go build -o /app/bin

# Expose a port
EXPOSE 8080

# Command to run the executable
CMD ["/app/bin/BUPTreasure"]
