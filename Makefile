.PHONY: build run dev clean docker-build docker-run deploy test

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

BINARY_NAME=jw-blog
DOCKER_IMAGE=jw-blog

# Build the application
build:
	CGO_ENABLED=1 $(GOBUILD) -o $(BINARY_NAME) ./cmd/server

# Run the application locally
run: build
	DB_PATH=./data/blog.db PORT=8080 ./$(BINARY_NAME)

# Development mode with auto-reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf ./data
	docker system prune -f

# Test
test:
	$(GOTEST) -v ./...

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Docker build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Docker run
docker-run: docker-build
	mkdir -p ./data
	docker run --rm -p 8080:8080 -v $(PWD)/data:/data $(DOCKER_IMAGE)

# Docker compose
docker-dev:
	mkdir -p ./data
	docker-compose up --build

# Initialize database with sample data
init-db:
	mkdir -p ./data
	@echo "Creating sample blog posts..."
	@sqlite3 ./data/blog.db <<< "$(shell cat scripts/sample_data.sql)"

# Deploy to production (assumes you have SSH access to your server)
deploy:
	@echo "Building production image..."
	docker build -t $(DOCKER_IMAGE):latest .
	@echo "Saving image to tar..."
	docker save $(DOCKER_IMAGE):latest | gzip > $(DOCKER_IMAGE).tar.gz
	@echo "Copying to server..."
	scp -i ~/.ssh/jw-web-key.pem $(DOCKER_IMAGE).tar.gz deploy.sh ec2-user@ssh.jared-wallace.com:~/
	@echo "Deploying on server..."
	ssh -i ~/.ssh/jw-web-key.pem ec2-user@ssh.jared-wallace.com 'docker load < $(DOCKER_IMAGE).tar.gz && chmod +x deploy.sh && ./deploy.sh'
	@echo "Cleaning up..."
	rm -f $(DOCKER_IMAGE).tar.gz

# Optimize images (requires imageoptim-cli or similar)
optimize-images:
	@echo "Optimizing images..."
	@find static/images -name "*.jpg" -o -name "*.png" | xargs -I {} sh -c 'cwebp -q 80 {} -o {}.webp'

# Create directory structure
setup:
	mkdir -p cmd/server
	mkdir -p templates
	mkdir -p static/{css,js,images}
	mkdir -p data
	mkdir -p scripts
	@echo "Directory structure created!"

# Development helpers
install-tools:
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run
