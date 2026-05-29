.PHONY: test build dev

test:
	cd backend && go test ./...
	cd web && pnpm test

build:
	cd backend && go build ./cmd/server
	cd web && pnpm build

dev:
	./scripts/dev.sh
