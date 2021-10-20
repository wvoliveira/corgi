.PHONY: build-nextjs
build-nextjs:
	cd ui; \
	yarn install; \
	NEXT_TELEMETRY_DISABLED=1 yarn run export

.PHONY: build
build: build-nextjs
	go build .

.PHONY: urls
urls:
	go run ./urls/cmd/urls/main.go -http.addr :8080

urls-swag:
	swag init --parseInternal --dir "./" -g ./urls/cmd/urls/main.go -o urls/docs
