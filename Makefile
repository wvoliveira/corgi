.PHONY: build-nextjs
build-nextjs:
	cd ui; \
	yarn install; \
	NEXT_TELEMETRY_DISABLED=1 yarn run export

.PHONY: build
build: build-nextjs
	go build .

.PHONY: url
auth-pwd:
	go run ./api/v1/auth/pwd/cmd/pwd/main.go -http.addr :8080

auth-pwd-swag:
	swag init --parseInternal --dir "./" -g ./api/v1/auth/pwd/cmd/pwd/main.go -o api/v1/auth/pwd/docs

.PHONY: url
url:
	go run ./api/v1/url/cmd/url/main.go -http.addr :8080

url-swag:
	swag init --parseInternal --dir "./" -g ./api/v1/url/cmd/url/main.go -o api/v1/url/docs

.PHONY: user
user:
	go run ./api/v1/url/cmd/url/main.go -http.addr :8080

user-swag:
	swag init --parseInternal --dir "./" -g ./api/v1/user/cmd/user/main.go -o api/v1/user/docs

main-swag:
	swag init --parseInternal --dir "./" -g main.go -o docs

# Unit test API.
go-test:
	go test -v ./server
