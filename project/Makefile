AUTH_BINARY=authApp
USER_BINARY=userApp
ADMIN_BINARY=adminApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_auth build_user					
	@echo "Stopping docker images (if running...)"								
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Stopped containers succesfully"

## down_rm: stop and remove all of data about containers
down_rm:
	@echo "Stopping and remove all container...."
	docker-compose down -v
	@echo "Stoping and remove running containers"

## build_auth: builds the auth binary as a linux executable
build_auth:
	@echo "Building auth binary..."
	cd ../auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api
	@echo "Done!"

## build_user: builds the user binary as a linux executable
build_user:
	@echo "Building user binary..."
	cd ../user-service && env GOOS=linux CGO_ENABLED=0 go build -o ${USER_BINARY} ./cmd/api
	@echo "Done!"	

## build_user: builds the user binary as a linux executable
build_admin:
	@echo "Building user binary..."
	cd ../admin-service && env GOOS=linux CGO_ENABLED=0 go build -o ${ADMIN_BINARY} ./cmd/api
	@echo "Done!"	
