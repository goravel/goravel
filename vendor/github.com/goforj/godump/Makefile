.PHONY: modernize modernize-fix modernize-check

MODERNIZE_CMD = go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@v0.18.1

modernize: modernize-fix

modernize-fix:
	@echo "Running gopls modernize with -fix..."
	$(MODERNIZE_CMD) -test -fix ./...

modernize-check:
	@echo "Checking if code needs modernization..."
	$(MODERNIZE_CMD) -test ./...
