#!/bin/bash

set -euo pipefail

VERSION=$(git rev-parse HEAD)
IMAGE="ghcr.io/hpcsc/book-stocker-api:${VERSION}"

docker build \
  -t "${IMAGE}" \
  --build-arg \
  VERSION="${VERSION}" \
  .

if [ "${PUSH:-}" = "true" ]; then
  docker push "${IMAGE}"
fi
