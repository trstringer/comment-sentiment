LANGUAGE_KEY_FILE=./tests/language_key

include .env

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
		--language-endpoint $(LANGUAGE_ENDPOINT) \
		--app-id 5 \
		--app-keyfile $(LANGUAGE_KEY_FILE)

.PHONY: debug
debug:
	dlv debug . -- \
		--language-keyfile $(LANGUAGE_KEY_FILE) \
		--language-endpoint $(LANGUAGE_ENDPOINT) \
		--app-id 5 \
		--app-keyfile $(LANGUAGE_KEY_FILE)

.PHONY: clean
clean:
	rm -rf ./dist
