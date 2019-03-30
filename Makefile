.DEFAULT_GOAL := help

.PHONY: ci tools clean test do-cover cover build image help

NAME    = main
VERSION = 1.0.0
GOTOOLS = \
	github.com/golang/dep/cmd/dep \
	golang.org/x/tools/cmd/cover


create-kube-config: ## Remove old binary
	mkdir ~/.kube || true && ./create-k8s-config.sh

install-kubectl: ## Remove old binary
	curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/darwin/amd64/kubectl
	chmod +x ./kubectl
	sudo mv ./kubectl /usr/local/bin/kubectl

blue-green: ## Performs a blue greeen deployment
	./deployer.sh --set-deployment production

ci: clean test build docker install-kubectl create-kube-config blue-green ## Continous Integration Steps

clean: ## Remove old binary
	-@rm -f $(NAME); \
	find vendor/* -maxdepth 0 -type d -exec rm -rf '{}' \;

test:  ## Execute tests
	go test $$(go list ./... | grep -v vendor)

do-cover: ## Test report
	./scripts/cover.sh

cover: env do-cover env-stop ## [env do-cover env-stop]

build: clean test ## [clean test] Build binary file
	docker build -t thiagotr/fortune-backend .

docker: ## Build Docker image
	docker login -u ${DOCKER_LOGIN} -p ${DOCKER_PASSWORD}
	docker push thiagotr/fortune-backend

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'