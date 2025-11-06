.PHONY: test install cover lint bench generate install-termux clean update-deps

# Tool definitions
TOOLS := epub-lint song-converter cat-ascii magnum jp-proc versy
GENERATE_TOOLS := epub-lint jp-proc magnum song-converter

# Enhanced LDFLAGS for size reduction
LDFLAGS := -ldflags="-s -w"
BUILDFLAGS := -trimpath
GCFLAGS := -gcflags="-l=4"

BUILD_CMD = CGO_ENABLED=0 GOOS=linux go build $(BUILDFLAGS) $(LDFLAGS) $(GCFLAGS)

# Bash completion directory with fallback
BASH_COMPLETION_DIR := $(or $(BASH_COMPLETION_USER_DIR),$(HOME)/.bash_completion.d)

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
	@echo "Building go tools and generating bash completion"
	@echo "Using bash completion directory: $(BASH_COMPLETION_DIR)"
	@mkdir -p "$(BASH_COMPLETION_DIR)"
	@for tool in $(TOOLS); do \
		echo "Building $$tool..."; \
		$(BUILD_CMD) -o "$(HOME)/.local/bin/$$tool" ./$$tool/main.go || { \
			echo "Error: Failed to build $$tool"; \
			exit 1; \
		}; \
		echo "Generating completion for $$tool..."; \
		"$$tool" completion bash > "$(BASH_COMPLETION_DIR)/$$tool-completion" || { \
			echo "Warning: Failed to generate completion for $$tool"; \
		}; \
	done

	@echo ""
	@echo "Tools installed successfully"

generate:
	@echo "Generating documentation..."
	@for tool in $(GENERATE_TOOLS); do \
		echo "Generating docs for $$tool..."; \
		go run --tags="generate_doc" ./$$tool/main.go generate -g ./$$tool/ || { \
			echo "Error: Failed to generate docs for $$tool"; \
			exit 1; \
		}; \
	done

install-termux:
	@echo "Building epub-lint for Termux"
	@CGO_ENABLED=0 go build $(BUILDFLAGS) $(LDFLAGS) $(GCFLAGS) -o "${PREFIX}/bin/epub-lint" ./epub-lint/main.go
	@mkdir -p ${PREFIX}/share/bash-completion/completions
	@echo "Generating bash completion for epub-lint"
	@epub-lint completion bash > "${PREFIX}/share/bash-completion/completions/epub-lint"

clean:
	@echo "Cleaning built binaries..."
	@for tool in $(TOOLS); do \
		if [ -f "$(HOME)/.local/bin/$$tool" ]; then \
			echo "Removing $(HOME)/.local/bin/$$tool"; \
			rm -f "$(HOME)/.local/bin/$$tool"; \
		fi; \
	done
	@echo "Cleaning bash completions..."
	@for tool in $(TOOLS); do \
		if [ -f "$(BASH_COMPLETION_DIR)/$$tool-completion" ]; then \
			echo "Removing $(BASH_COMPLETION_DIR)/$$tool-completion"; \
			rm -f "$(BASH_COMPLETION_DIR)/$$tool-completion"; \
		fi; \
	done
	@if [ -f "$$PREFIX/bin/epub-lint" ]; then \
		echo "Removing $$PREFIX/bin/epub-lint"; \
		rm -f "$$PREFIX/bin/epub-lint"; \
	fi
	@echo "Cleanup complete"

update-deps:
	@echo "Installing all dependency updates"
	@go get -u ./...
	@go mod tidy
	@echo "Update complete"

golden:
	@echo "Generating golden files"
	@go run -tags "generate_golden" ./magnum/ golden -o ./magnum/internal
	@echo "Golden files generated"
