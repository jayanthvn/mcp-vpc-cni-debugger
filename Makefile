
APP_NAME = mcp-server
IMAGE_NAME = your-repo/$(APP_NAME)

export GOPROXY = direct

build:
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) ./cmd/mcp-server

docker-build:
	docker build -t $(IMAGE_NAME):latest .

docker-push:
	docker push $(IMAGE_NAME):latest

run:
	go run ./cmd/mcp-server
