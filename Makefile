.PHONY: build-nextjs
build-nextjs:
	cd ui; \
	yarn install; \
	NEXT_TELEMETRY_DISABLED=1 yarn run export

.PHONY: build
build: build-nextjs
	go build .

.PHONY: run
run:
	go run main.go
