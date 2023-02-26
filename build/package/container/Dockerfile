FROM node:14.21.3 AS ui

ENV NODE_OPTIONS --max_old_space_size=10240

WORKDIR /src/web
COPY . /src

RUN npm install
RUN NEXT_TELEMETRY_DISABLED=1 npm run export


FROM golang:1.20 AS binary

ARG VERSION
ARG REVISION
ARG DATE

ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod /src
COPY go.sum /src

RUN go mod download

COPY --from=ui /src/ /src

RUN go build -o /opt/corgi/corgi /src/cmd/corgi/*.go


FROM alpine:3.13

LABEL name="corgi"
LABEL description="A URL shortener app."
LABEL version=$VERSION
LABEL revision=$REVISION
LABEL ksdate=$DATE

RUN /sbin/apk --no-cache add ca-certificates curl iproute2 tini
RUN addgroup -g 1000 application && adduser -u 1000 -SDH -G application application

COPY --from=binary /opt/corgi/corgi /opt/corgi/corgi

WORKDIR /opt/corgi/

USER application
EXPOSE 8081

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["./corgi"]

HEALTHCHECK --interval=5m --timeout=3s \
  CMD curl -f http://localhost:8081/ || exit 1
