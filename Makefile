REGISTRY?=registry.gitlab.com/pardacho/secure-portal
APP_VERSION?=latest
BUILD?=go build -ldflags="-w -s"

.PHONY: proto

default: build

build: build-auth

build-auth: format lint
	$(BUILD) -o auth-server cmd/auth/main.go

proto:
	protoc -I proto services.proto --go_out=plugins=grpc:services

deps:
	go mod download

test:
	go test $$(go list ./... | grep -v /vendor/)

format:
	go fmt $$(go list ./... | grep -v /vendor/)

vet:
	go vet $$(go list ./... | grep -v /vendor/)

lint:
	golint -set_exit_status -min_confidence 0.3 $$(go list ./... | grep -v /vendor/)

registry: registry-build registry-push

registry-build:
	docker build --pull -f docker/auth/Dockerfile -t $(REGISTRY)/auth:$(APP_VERSION) .

registry-pull:
	docker pull $(REGISTRY)/auth:$(APP_VERSION)

registry-push:
	docker push $(REGISTRY)/auth:$(APP_VERSION)

registry-clear:
	docker image rm -f $(REGISTRY)/auth:$(APP_VERSION)

stop:
	docker-compose stop

stop-prod:
	docker stack rm app

dev:
	docker-compose build
	docker-compose up -d
	clear
	@echo ""
	@echo "starting command line:"
	@echo "** when finish exist and run: make stop**"
	@echo ""
	docker-compose exec server sh

prod:
	docker stack deploy --compose-file docker-stack.yml secure-portal --with-registry-auth
	clear
	@echo ""
	@echo "commands:"
	@echo "- make stop-prod"
	@echo ""