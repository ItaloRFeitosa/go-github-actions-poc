# Get the host user and group IDs
HOST_USER_ID := $(shell id -u)
HOST_GROUP_ID := $(shell id -g)
COMPOSE_USER :=  $(HOST_USER_ID):$(HOST_GROUP_ID)
LOCAL_COMPOSE_FILE := ./deployments/docker-compose.local.yml
CI_COMPOSE_FILE := ./deployments/docker-compose.ci.yml
.PHONY: destroy e2e build dev ci

build:
	go build -o ./out/app ./cmd/app/main.go

dev:
	DOCKER_BUILDKIT=0 COMPOSE_USER=$(COMPOSE_USER) docker compose --project-directory . -f $(LOCAL_COMPOSE_FILE) up -d

dev_stop:
	DOCKER_BUILDKIT=0 COMPOSE_USER=$(COMPOSE_USER) docker compose --project-directory . -f $(LOCAL_COMPOSE_FILE) stop

dev_logs:
	DOCKER_BUILDKIT=0 COMPOSE_USER=$(COMPOSE_USER) docker compose --project-directory . -f $(LOCAL_COMPOSE_FILE) logs app --follow

dev_build:
	DOCKER_BUILDKIT=0 COMPOSE_USER=$(COMPOSE_USER) docker compose --project-directory . -f $(LOCAL_COMPOSE_FILE) up -d --build

ci:
	DOCKER_BUILDKIT=0 docker compose --project-directory . -p promos_app_ci -f $(CI_COMPOSE_FILE) up --attach test --build --abort-on-container-exit --exit-code-from test

ci_nobuild:
	DOCKER_BUILDKIT=0 docker compose --project-directory . -p promos_app_ci -f $(CI_COMPOSE_FILE) up --attach test --abort-on-container-exit --exit-code-from test

ci_test_logs:
	DOCKER_BUILDKIT=0 docker compose --project-directory . -p promos_app_ci -f $(CI_COMPOSE_FILE) logs test

ci_down:
	DOCKER_BUILDKIT=0 docker compose --project-directory . -p promos_app_ci -f $(CI_COMPOSE_FILE) down -v

destroy:
	docker-compose --project-directory . -f $(LOCAL_COMPOSE_FILE) down
	sudo rm -rf ./.docker

e2e:
	go test ./test/e2e/... -v -count=1