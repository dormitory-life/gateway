IMAGE ?= gateway-svc
TAG ?= local

.PHONY: local local-build local-run build run start

start: build run

build:
	@echo "Building gateway svc..."
	@mkdir -p .bin
	@cd $(CURDIR) && go build -o .bin/main cmd/main.go

run:
	@echo "Starting gateway svc..."
	@cd $(CURDIR) && go run cmd/main.go configs/config.yaml

local: local-build local-run

local-build:
	@echo "Building gateway svc..."
	@cd $(CURDIR) && go build -o .bin/main cmd/main.go

local-run:
	@echo "Starting gateway svc..."
	@cd $(CURDIR) && go run cmd/main.go configs/local.yaml

local-stop:
	@echo "Stopping gateway svc..."
	@if lsof -ti:$$8080 >/dev/null 2>&1; then \
			echo "Порт $$8080 занят. Убиваю..."; \
			kill -9 $$(lsof -ti:$$8080) 2>/dev/null || true; \
		else \
			echo "Порт $$8080 свободен"; \
		fi; \


.PHONY: docker-build
docker-build:
	@echo "Building docker image..."
	@docker build -t $(IMAGE):$(TAG) .