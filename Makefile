.DEFAULT_GOAL := help

.PHONY: ci tools clean test do-cover cover build image help

NAME    = main
VERSION = 1.0.0



clean-local: ## cleans local test infra
	docker rm -f mongo-tests

local-test: clean-local ## Run all tests locally using docker
	docker run --name mongo-tests --network=host -d mongo
	./scripts/cover-script.sh
	docker rm -f mongo-tests

create-kube-config: ## Remove old binary
	mkdir ~/.kube || true && ./scripts/create-k8s-config.sh

install-kubectl: ## Remove old binary
	curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.14.0/bin/linux/amd64/kubectl
	chmod +x ./kubectl
	sudo mv ./kubectl /usr/local/bin/kubectl

blue-green: ## Performs a blue greeen deployment
	./scripts/deployer.sh --set-deployment production

ci: clean test build docker install-kubectl create-kube-config blue-green ## Continous Integration Steps

clean: ## Remove old binary
	-@rm -f $(NAME); \
	find vendor/* -maxdepth 0 -type d -exec rm -rf '{}' \;

test: do-cover ## Execute tests
	go test -cover $$(go list ./... | grep -v vendor)

build: clean ## [clean test] Build binary file
	docker build -t thiagotr/fortune-backend .

docker-login: ## [clean test] Build binary file
	docker login -u ${DOCKER_LOGIN} -p ${DOCKER_PASSWORD}

docker-push: build docker-login ## Build Docker image
	docker push thiagotr/fortune-backend

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'