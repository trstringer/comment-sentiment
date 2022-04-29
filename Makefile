LANGUAGE_KEY_FILE=./test/language-key-file.txt

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
run: generate
	go run . --language-key $(LANGUAGE_KEY_FILE)

.PHONY: clean
clean:
	rmdir ./dist
