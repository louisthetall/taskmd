.PHONY: install-dev install-dev-full check check-lite test lint sync-spec docker-build docker-run

# Install development binary (delegates to apps/cli)
install-dev:
	$(MAKE) -C apps/cli install-dev

# Install development binary with embedded web assets
install-dev-full:
	$(MAKE) -C apps/cli install-dev-full

# Quick check: compile and test CLI + SDK
check-lite:
	cd apps/cli && go build ./...
	cd apps/cli && go test ./...
	cd sdk/go && go build ./...
	cd sdk/go && go test ./...

# Run all checks (CLI tests, lint, vet + SDK + web tests + docs build + Docker build)
check: check-lite
	$(MAKE) -C apps/cli lint
	cd apps/web && npx vitest run
	cd apps/docs && pnpm build
	docker build -t taskmd:ci-check .

# Run tests only
test:
	$(MAKE) -C apps/cli test
	cd sdk/go && go test ./...

# Run linter only
lint:
	$(MAKE) -C apps/cli lint

# Sync spec copies from docs/taskmd_specification.md
sync-spec:
	$(MAKE) -C apps/cli sync-spec

# Build Docker image
docker-build:
	docker build -t taskmd:local .

# Run Docker container (mount ./tasks as read-only)
docker-run: docker-build
	docker run --rm -p 8080:8080 -v ./tasks:/tasks:ro taskmd:local
