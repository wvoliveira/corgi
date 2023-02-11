FROM node:16.14.0 AS ui

WORKDIR /workspace/web

COPY . /workspace

ENV NODE_OPTIONS --max_old_space_size=10240

RUN yarn install
RUN NEXT_TELEMETRY_DISABLED=1 yarn run export


FROM golang:1.18.5 AS binary

ARG commit
ARG version
ARG ksdate

ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
ENV CGO_ENABLED=0

COPY --from=ui "/workspace" /workspace
WORKDIR /workspace

RUN rm -rfv /workspace/cmd/web; mv /workspace/web/dist /workspace/cmd/corgi/web
RUN go build -o /app /workspace/cmd/corgi/*.go


FROM alpine:3.13

ARG version

LABEL version=$version
LABEL ksdate=$ksdate
LABEL description="Corgi. A shortener app."

RUN /sbin/apk --no-cache add ca-certificates curl iproute2 tini

COPY --from=binary "/app" /
COPY --from=binary "/workspace/rbac_*" /

WORKDIR /

EXPOSE 8081

ENTRYPOINT ["/sbin/tini", "--"]

CMD ["/app"]

HEALTHCHECK --interval=5m --timeout=3s \
  CMD curl -f http://localhost:8081/ || exit 1
