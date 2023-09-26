build:
	@echo "Build chat server..."
	docker build -t chat-server . -f ./server/Dockerfile

up:
	@echo "Run chat server..."
	docker run -d --name chat-server -p 8080:8080 chat-server

up_build: build up

down:
	@echo "Stop chat server..."
	docker stop chat-server
	docker rm chat-server

watch:
	@echo "Watch chat server..."
	docker logs -f chat-server