FROM docker.io/library/alpine:3.19 as runtime

RUN \
  apk add --update --no-cache \
    bash \
    coreutils \
    curl \
    ca-certificates \
    tzdata

ENTRYPOINT ["promtool-jsonnet"]
COPY promtool-jsonnet /usr/bin/

COPY .cache/prometheus /usr/lib/prometheus
ENV PJ_PROMTOOL_PATH="/usr/lib/prometheus/promtool"

COPY jsonnet /usr/lib/jsonnet
ENV PJ_JSONNET_PATH="/usr/lib/jsonnet"

USER 65536:0
