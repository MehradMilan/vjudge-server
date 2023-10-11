FROM golang:1.21

# # Set the URL for the repository
# ENV VJUDGE_CORE_URL="https://github.com/sorousherafat/vjudge-core"

# # Clone the repository
# RUN git clone $VJUDGE_CORE_URL /core

# Set working directory
WORKDIR /app

# Copy the local source files to the container
COPY . .

# Download dependencies and tidy up
RUN go mod download
RUN go mod tidy

# Build the application
RUN go build -o webhook ./cmd/webhook

# Expose the application's port
EXPOSE 8000

# Command to run the application
CMD ["./webhook", "config/config-webhook"]
