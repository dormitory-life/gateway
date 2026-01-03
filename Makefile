.PHONY: local local-build local-run build run start

start: build run

build:
	@echo "Building gateway svc..."
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