.PHONY: docker-up
docker-up:
	docker-compose -f docker-compose.yaml up --build

.PHONY: docker-down
docker-down: ## Stop docker containers and clear artefacts.
	docker-compose -f docker-compose.yaml down
	docker system prune 

.PHONY: build-up
build-up:
	GOARCH=amd64 GOOS=linux go build -tags=jsoniter -o ./build/main cmd/web/*.go 
	