IMAGE_NAME=vjudge
IMAGE_TAG=v1.0.0

run:
	@echo "Running the Go application..."
	@go run cmd/webhook

build-image:
	@echo "Building the Docker image..."
	@sudo docker build . -t $(IMAGE_NAME):$(IMAGE_TAG)

compose-up:
	@echo "Starting the application using Docker Compose..."
	@sudo docker-compose up

.PHONY: run build-image compose-up
