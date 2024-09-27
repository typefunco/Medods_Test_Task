start_server:
	go run authService/cmd/main.go

docker-build:
	@echo "Building Docker image..."
	docker build -t auth-service -f DockerFile .

docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 auth-service


