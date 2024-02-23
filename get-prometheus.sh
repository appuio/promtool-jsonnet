#!/bin/bash

PROMETHEUS_VERSION="2.40.7"
PROMETHEUS_DIST="`go env GOOS`"
PROMETHEUS_ARCH="`go env GOARCH`"
PROMETHEUS_DOWNLOAD_LINK="https://github.com/prometheus/prometheus/releases/download/v${PROMETHEUS_VERSION}/prometheus-${PROMETHEUS_VERSION}.${PROMETHEUS_DIST}-${PROMETHEUS_ARCH}.tar.gz"

CACHE_DIR="${1:-.cache}"

mkdir -p "${CACHE_DIR}"
curl -fsSLo "${CACHE_DIR}"/prometheus.tar.gz ${PROMETHEUS_DOWNLOAD_LINK}
tar -xzf "${CACHE_DIR}"/prometheus.tar.gz -C .cache
mv "${CACHE_DIR}"/prometheus-${PROMETHEUS_VERSION}.${PROMETHEUS_DIST}-${PROMETHEUS_ARCH} .cache/prometheus
rm -rf "${CACHE_DIR}"/*.tar.gz
