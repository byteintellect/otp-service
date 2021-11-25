PROJECT_ROOT            ?= $(PWD)
CURDIR ?= $(shell pwd)
BUILD_PATH              ?= bin
DOCKERFILE_PATH         ?= $(CURDIR)

# configuration for image names
USERNAME                ?= $(USER)
GIT_COMMIT              ?= $(shell git describe --always --tags || echo pre-commit)
IMAGE_VERSION           ?= $(GIT_COMMIT)

# configuration for server binary and image
SERVER_BINARY           ?= $(BUILD_PATH)/server
SERVER_PATH             ?= $(PROJECT_ROOT)/cmd/server
SERVER_DOCKERFILE       ?= $(DOCKERFILE_PATH)/Dockerfile
SERVER_IMAGE            := user_service
DOCKER_REGISTRY_RUL 	:= docker-registry.cashfree.com

.PHONY up_migration:
up_migration:
	MIGRATIONS_DIR="$(MIGRATIONS_DIR)" DATABASE_NAME="$(DATABASE_NAME)" DATABASE_DSN="$(DATABASE_DSN)" dbmigrator migrate up

.PHONY down_migration:
down_migration:
	MIGRATIONS_DIR="$(MIGRATIONS_DIR)" DATABASE_NAME="$(DATABASE_NAME)" DATABASE_DSN="$(DATABASE_DSN)" dbmigrator migrate down

.PHONY docker_build:
docker_build:
	@docker build --build-arg GITHUB_USER="$(GITHUB_USER)" --build-arg GITHUB_TOKEN="$(GITHUB_TOKEN)" -t "${DOCKER_REGISTRY_URL}/${SERVER_IMAGE}":"${IMAGE_VERSION}" .

.PHONY docker_login:
docker_login:
	@docker login ${DOCKER_REGISTRY_URL} -u ${REGISTRY_USERNAME} -p ${REGISTRY_PASSWORD}

.PHONY docker_push:
docker_push:
	@docker push ${DOCKER_REGISTRY_URL}/${SERVER_IMAGE} --all-tags

.PHONY docker_login_and_push:
docker_login_and_push: docker_login docker_push

.PHONY docker_build_and_push:
docker_build_and_push: docker_build docker_login docker_push

.PHONY test:
test:
	@go test ./... || (echo "go test failed $$?"; exit 1)

.PHONY lint:
lint:
	@go fmt ./... && go vet ./... || (echo "go vet failed $$?"; exit 1)

.PHONY swagger-ui:
swagger-ui:
	swagger generate spec -m -i swagger/swagger.yml -o ./swagger.json

replace_version:
	. ./script.sh; replace helm/user-svc

commit:
	echo ${LATEST_VERSION}
	git config --global user.email "kumar.varalakshmi@outlook.com"
	git config --global user.name "Kumar D"
	git add . -A
	git commit -m "next release to $(shell git describe --always --tags || echo pre-commit)"
	git push origin HEAD

upgrade:
	helm upgrade --install -f helm/user-svc/values.yaml user-svc ./helm/user-svc --namespace ${NAMESPACE}

replace_and_commit: replace_version upgrade commit
