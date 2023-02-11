export CGO_ENABLED := 0
export NEXT_TELEMETRY_DISABLED := 1


.PHONY: build
build:
	go build -ldflags "-s -w" -o corgi ./cmd/corgi/*.go


.PHONY: build-web
build-web:
	cd web && \
	npm install --frozen-lockfile && \
	npm run export && \
	rm -rfv ../cmd/corgi/web; \
    mv dist ../cmd/corgi/web


.PHONY: clean
clean:
	rm -f spitz
	rm -rf ./cmd/corgi/web
	rm -rf ./web/dist
	rm -rf ./web/.next


.PHONY: load-test-2m
load-test-2m:
	k6 run ./scripts/load-test/2-min.js


# First run some requests and run pprof-graph to get the numbers in svg graph.
# Good to debug your application to get where it's speding most of time.
# Oh god, my english.. 
.PHONY: pprof-graph
pprof-graph:
	pprof -web "http://:8081/debug/pprof/profile?seconds=5"


local-dep:
	docker-compose -f deployments/container/docker-compose.yaml up
