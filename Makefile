
.PHONY: test install cover install-deps lint

PLAYWRIGHT_VERSION=$(shell go list -m github.com/playwright-community/playwright-go | awk '{print $$2}')

test:
	go test ./... -tags "unit"

# this is just meant to give an idea whether or not something has tests in it.
# It is not meant to be used for 100% test coverage. Some folders will be better tested than others.
cover:
	go test -cover ./... -tags "unit"

lint:
	golangci-lint run ./...

install:
	@echo "Building go tools"
	@go build -o "${HOME}/.local/bin/ebook-lint" ./ebook-lint/main.go
	@go build -o "${HOME}/.local/bin/git-helper" ./git-helper/main.go
	@go build -o "${HOME}/.local/bin/song-converter" ./song-converter/main.go
	@go build -o "${HOME}/.local/bin/cat-ascii" ./cat-ascii/main.go
	@go build -o "${HOME}/.local/bin/magnum" ./magnum/main.go
	@go build -o "${HOME}/.local/bin/jp-proc" ./jp-proc/main.go
	
	@mkdir -p ${BASH_COMPLETION_USER_DIR}

	@echo "Generating the bash completion for the tools"
	@ebook-lint completion bash > "${BASH_COMPLETION_USER_DIR}/ebook-lint-completion"
	@git-helper completion bash > "${BASH_COMPLETION_USER_DIR}/git-helper-completion"
	@song-converter completion bash > "${BASH_COMPLETION_USER_DIR}/song-converter-completion"
	@cat-ascii completion bash > "${BASH_COMPLETION_USER_DIR}/cat-ascii-completion"
	@magnum completion bash > "${BASH_COMPLETION_USER_DIR}/magnum-completion"
	@jp-proc completion bash > "${BASH_COMPLETION_USER_DIR}/jp-proc-completion"

install-deps:
	@go run github.com/playwright-community/playwright-go/cmd/playwright@${PLAYWRIGHT_VERSION} install --with-deps
