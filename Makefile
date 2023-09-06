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

.PHONY: build-up-prod
build-up-prod:
	GOARCH=amd64 GOOS=linux go build -tags=jsoniter -o ./build/prod/main cmd/web/*.go 

.PHONY: build-up-dev
build-up-dev:
	GOARCH=amd64 GOOS=linux go build -tags=jsoniter -o ./build/dev/main cmd/web/*.go 

.PHONY: build-up-test
build-up-test:
	GOARCH=amd64 GOOS=linux go build -tags=jsoniter -o ./build/test/main cmd/web/*.go 
	

