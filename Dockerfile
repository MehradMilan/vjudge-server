FROM golang:1.23

# Install Icarus Verilog
RUN apt -y update && apt -y install iverilog

WORKDIR /lib
RUN git clone https://github.com/sorousherafat/libvcd.git && make -C libvcd && make -C libvcd install
RUN git clone https://github.com/sorousherafat/libvjudge.git && make -C libvjudge && make -C libvjudge install

# Clone the repository
WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target="/go/pkg/mod" go mod download

COPY . .
RUN --mount=type=cache,target="/go/pkg/mod" go build -o webhook ./cmd/webhook

# Expose the application's port
EXPOSE 8000

# Command to run the application
CMD ["./webhook", "config/config-webhook.json"]
