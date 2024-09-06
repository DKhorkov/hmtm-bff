all: app

app:
	@make -s help

help:
	@echo "make help                         - show this help."
	@echo "make build                        - build docker and it's environment."
	@echo "make run							 - run docker."
	@echo "make clean                        - cleanup docker kitchen."
	@echo "make network                      - create docker network for project."
	@echo "make prod                         - run docker for production purposes."

network:
	sudo docker network create hmtmbff_network

build:
	sudo docker compose -f ../build/package/docker-compose.yml build

run:
	sudo docker compose -f ../build/package/docker-compose.yml up

clean:
	sudo docker system prune -a --volumes

prod: build run
