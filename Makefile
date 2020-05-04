gen:
	@echo "Running go generate..."
	go generate github.com/jmgilman/gcli/vault/auth

test:
	@echo "Running all tests..."
	go test ./vault/auth/... ./vault/client/...
