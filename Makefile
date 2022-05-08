LANGUAGE_KEY_FILE=./tests/language_key
LANGUAGE_ENDPOINT=$(shell cat ./tests/language_endpoint)
INFRA_OUT_PLAN_FILE=./infra.out

.PHONY: build
build: generate
	mkdir -p ./dist
	go build -o ./dist/comment-sentiment .

.PHONY: test
test:
	go test -v ./...

.PHONY: generate
generate:
	./scripts/generate.sh

.PHONY: run
run:
	go run . \
		--language-keyfile $(LANGUAGE_KEY_FILE) \
		--language-endpoint $(shell terraform -chdir=./infra output -raw language_endpoint) \
		--app-id 5 \
		--app-keyfile $(LANGUAGE_KEY_FILE)

.PHONY: debug
debug:
	dlv debug . -- \
		--language-keyfile $(LANGUAGE_KEY_FILE) \
		--language-endpoint $(shell terraform -chdir=./infra output -raw language_endpoint) \
		--app-id 5 \
		--app-keyfile $(LANGUAGE_KEY_FILE)

.PHONY: clean
clean:
	rm -rf ./dist

.PHONY: infra-plan
infra-plan:
	terraform -chdir=./infra plan

.PHONY: infra-apply
infra-apply:
	./scripts/infra_apply.sh

.PHONY: infra-clean
infra-clean:
	./scripts/infra_clean.sh

.PHONY: image-build
image-build: build
	docker build -t \
		$$(terraform -chdir=./infra output -raw acr_endpoint)/comment-sentiment:$$(./dist/comment-sentiment -v) .

.PHONY: image-push
image-push: image-build
	ACR=$$(terraform -chdir=./infra output -raw acr_endpoint) && \
	az acr login -n $$ACR && \
	docker push $$ACR/comment-sentiment:$$(./dist/comment-sentiment -v)

.PHONY: deploy
deploy:
	./scripts/deploy.sh

.PHONY: chart-lint
chart-lint:
	./scripts/chart_lint.sh

.PHONY: dns
dns:
	./scripts/dns.sh
