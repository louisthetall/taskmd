.PHONY: check test lint sync-spec

# Run all checks (tests, lint, vet)
check:
	$(MAKE) -C apps/cli check

# Run tests only
test:
	$(MAKE) -C apps/cli test

# Run linter only
lint:
	$(MAKE) -C apps/cli lint

# Sync spec copies from docs/taskmd_specification.md
sync-spec:
	$(MAKE) -C apps/cli sync-spec
