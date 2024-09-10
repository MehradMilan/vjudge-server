FROM golang:1.23

# Install Icarus Verilog
RUN apt -y update && apt -y install iverilog

WORKDIR /lib
RUN git clone https://github.com/sorousherafat/libvcd.git && cd libvcd && make && make install

WORKDIR /lib
RUN git clone https://github.com/sorousherafat/libvjudge.git && cd libvjudge && make && make install

# Clone the repository
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Copy the local source files to the container
RUN go build -o webhook ./cmd/webhook

# Expose the application's port
EXPOSE 8000

# Command to run the application
CMD ["./webhook", "config/config-webhook.json"]
