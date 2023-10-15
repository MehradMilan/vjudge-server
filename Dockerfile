FROM golang:1.21.1

# Install Icarus Verilog
RUN apt -y update && apt -y install iverilog

# Clone the repository
WORKDIR /app
COPY . .

# Set working directory
WORKDIR /app/lib/libvjudge
RUN make

WORKDIR /app/lib/libvcd
RUN make

WORKDIR /app/

# Copy the local source files to the container

RUN go build -mod=vendor -o webhook ./cmd/webhook

# Expose the application's port
EXPOSE 8000

# Command to run the application
CMD ["./webhook", "config/config-webhook.json"]
