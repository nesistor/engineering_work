AUTH_BINARY=authApp
USER_BINARY=userApp
ADMIN_BINARY=adminApp
MAIL_BINARY=mailApp
LOGGER_BINARY=loggerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_auth build_user	build_admin build_mail build_logger 			
	@echo "Stopping docker images (if running...)"								
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## build: create new version apps:
build: build_auth build_user build_admin build_mail build_logger
	@echo "Succesfully compiled apps from binary!"


## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Stopped containers succesfully"

## down_rm: stop and remove all of data about containers
down_rm:
	@echo "Stopping and remove all container...."
	docker compose down -v
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

## build_admin: builds the user binary as a linux executable
build_admin:
	@echo "Building admin binary..."
	cd ../admin-service && env GOOS=linux CGO_ENABLED=0 go build -o ${ADMIN_BINARY} ./cmd/api
	@echo "Done!"	

## build_mailer: builds the user binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Done!"	

## build_logger: builds the user binary as a linux executable
build_logger:
	@echo "Building logger binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BINARY} ./cmd/api
	@echo "Done!"	
