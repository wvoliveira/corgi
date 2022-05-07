export CGO_ENABLED = 1
export NEXT_TELEMETRY_DISABLED = 1

.PHONY: build
build: build-web
	go build ./cmd/corgi

.PHONY: build-web
build-web:
	cd web && \
	npm install --frozen-lockfile && \
	npm run export && \
    mv dist ../cmd/corgi/web

.PHONY: clean
clean:
	rm -f spitz
	rm -rf ./cmd/corgi/web
	rm -rf ./web/dist
	rm -rf ./web/.next
