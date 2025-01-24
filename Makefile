.PHONY: test install cover lint bench generate install-termux

test:
	go test ./... -tags "unit"

# this is just meant to give an idea whether or not something has tests in it.
# It is not meant to be used for 100% test coverage. Some folders will be better tested than others.
cover:
	go test -cover ./... -tags "unit"

lint:
	golangci-lint run ./...

bench:
	go test ./... -bench=. -tags="unit" -count=20 -run=^$

install:
	@echo "Building go tools"
	@go build -ldflags="-s -w" -o "${HOME}/.local/bin/ebook-lint" ./ebook-lint/main.go
	@go build -ldflags="-s -w" -o "${HOME}/.local/bin/git-helper" ./git-helper/main.go
	@go build -ldflags="-s -w" -o "${HOME}/.local/bin/song-converter" ./song-converter/main.go
	@go build -ldflags="-s -w" -o "${HOME}/.local/bin/cat-ascii" ./cat-ascii/main.go
	@go build -ldflags="-s -w" -o "${HOME}/.local/bin/magnum" ./magnum/main.go
	@go build -ldflags="-s -w" -o "${HOME}/.local/bin/jp-proc" ./jp-proc/main.go
	
	@mkdir -p ${BASH_COMPLETION_USER_DIR}

	@echo "Generating the bash completion for the tools"
	@ebook-lint completion bash > "${BASH_COMPLETION_USER_DIR}/ebook-lint-completion"
	@git-helper completion bash > "${BASH_COMPLETION_USER_DIR}/git-helper-completion"
	@song-converter completion bash > "${BASH_COMPLETION_USER_DIR}/song-converter-completion"
	@cat-ascii completion bash > "${BASH_COMPLETION_USER_DIR}/cat-ascii-completion"
	@magnum completion bash > "${BASH_COMPLETION_USER_DIR}/magnum-completion"
	@jp-proc completion bash > "${BASH_COMPLETION_USER_DIR}/jp-proc-completion"

generate:
	@go run --tags="generate_doc" ./ebook-lint/main.go generate -g ./ebook-lint/
	@go run --tags="generate_doc" ./jp-proc/main.go generate -g ./jp-proc/
	@go run --tags="generate_doc" ./magnum/main.go generate -g ./magnum/
	@go run --tags="generate_doc" ./song-converter/main.go generate -g ./song-converter/

install-termux:
	@echo "Building ebook-lint for Termux"
	@go build -o "${PREFIX}/bin/ebook-lint" ./ebook-lint/main.go
	@mkdir -p ${PREFIX}/share/bash-completion/completions
	@echo "Generating bash completion for ebook-lint"
	@ebook-lint completion bash > "${PREFIX}/share/bash-completion/completions/ebook-lint"
