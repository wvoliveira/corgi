FROM node:14.16.0-alpine3.11 AS ui-builder

COPY . /workspace
WORKDIR /workspace/ui

ENV NODE_OPTIONS --max_old_space_size=10240

# remove stale node modules and dist from (potential) previous builds
RUN yarn install
RUN NEXT_TELEMETRY_DISABLED=1 yarn run export

FROM golang:1.17.2-alpine3.14 AS binary-builder

ARG commit
ARG version
ARG ksdate

ENV GOOS linux
ENV GOARCH amd64
ENV GO111MODULE on
ENV GOBIN /go/bin

COPY --from=ui-builder "/workspace" /workspace
WORKDIR /workspace

RUN apk --no-cache add ca-certificates gcc musl-dev

RUN go build -o /go/build/server .

FROM alpine:3.13

ARG version

LABEL version=$version
LABEL ksdate=$ksdate
LABEL description="Redir server for ELGA inc."

RUN apk --no-cache add ca-certificates curl iproute2 tini

WORKDIR /server/
COPY --from=binary-builder "/go/build/server" /server/server
EXPOSE 8080

ENTRYPOINT ["/sbin/tini", "--", "/server/server"]

HEALTHCHECK --interval=5m --timeout=3s \
  CMD curl -f http://localhost:8080/ || exit 1