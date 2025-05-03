# Terraform Provider for Spotify Makefile

BINARY_NAME=terraform-provider-spotify
VERSION=$(shell grep -E "Version = \"[0-9\.]+\"" version/version.go | sed -E 's/.*Version = "([0-9\.]+)".*/\1/')
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
PLUGIN_DIR=$(HOME)/.terraform.d/plugins/local/spotify/$(VERSION)/$(OS_ARCH)
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_DATE=$(shell date -u '+%Y-%m-%d')

.PHONY: build clean install test testacc lint fmt deps release scheduled-refresh

default: build

build:
	go build -ldflags "-X github.com/ashrafxbilal/terraform-provider-spotify/version.GitCommit=$(GIT_COMMIT) -X github.com/ashrafxbilal/terraform-provider-spotify/version.BuildDate=$(BUILD_DATE)" -o $(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)

install: build
	mkdir -p $(PLUGIN_DIR)
	cp $(BINARY_NAME) $(PLUGIN_DIR)/

# Run unit tests
test:
	go test ./... -v

# Check test coverage (aim for >70%)
coverage:
	chmod +x ./scripts/check_coverage.sh
	./scripts/check_coverage.sh

# Run acceptance tests
testacc:
	TF_ACC=1 go test ./... -v -timeout 120m

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Check for outdated dependencies
deps:
	./scripts/update_dependencies.sh

# Validate version follows semantic versioning (MAJOR.MINOR.PATCH)
validate-version:
	@echo "Validating version: $(VERSION)"
	@if ! echo $(VERSION) | grep -qE '^[0-9]+(\.[0-9]+){2}(-(alpha|beta|rc)[0-9]*)?$$'; then \
		echo "Error: Version $(VERSION) does not follow semantic versioning (MAJOR.MINOR.PATCH)"; \
		echo "Please update version/version.go with a valid semantic version"; \
		exit 1; \
	fi
	@echo "Version $(VERSION) is valid"

# Release a new version
release: validate-version
	git tag v$(VERSION)
	git push origin v$(VERSION)

# Run the auth proxy
auth-proxy:
	cd spotify_auth_proxy && go run spotify-auth.go

# Run the example
example: install
	cd examples/basic_playlist && terraform init && terraform apply

# Run the scheduled playlist refresh
scheduled-refresh: install
	cd examples/scheduled_refresh && terraform init && terraform apply

# Help target
help:
	@echo "Available targets:"
	@echo "  build             - Build the provider"
	@echo "  clean             - Remove build artifacts"
	@echo "  install           - Install the provider to the local Terraform plugin directory"
	@echo "  test              - Run unit tests"
	@echo "  testacc           - Run acceptance tests"
	@echo "  lint              - Run linter"
	@echo "  fmt               - Format code"
	@echo "  deps              - Check for outdated dependencies"
	@echo "  release           - Tag and push a new release"
	@echo "  auth-proxy        - Run the Spotify auth proxy"
	@echo "  example           - Run the example playlist generator"
	@echo "  scheduled-refresh - Run the scheduled playlist refresh"