include .env

up:
	@echo "Starting up the container..."
	docker-compose up --build -d

down:
	@echo "Stopping the container..."
	docker-compose down

start-gateway:
	cd gateway && air

start-orders:
	cd orders && air