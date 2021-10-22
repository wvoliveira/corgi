.PHONY: build-nextjs
build-nextjs:
	cd ui; \
	yarn install; \
	NEXT_TELEMETRY_DISABLED=1 yarn run export

.PHONY: build
build: build-nextjs
	go build .

.PHONY: url
url:
	go run ./url/cmd/url/main.go -http.addr :8080

url-swag:
	swag init --parseInternal --dir "./" -g ./url/cmd/url/main.go -o url/docs

main-swag:
	swag init --parseInternal --dir "./" -g main.go -o docs