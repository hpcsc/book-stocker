#!/bin/bash

set -euo pipefail

COMPONENT_NAME=${1}
shift
BUILD_ARGS="$@"

VERSION=$(git rev-parse HEAD)
IMAGE="ghcr.io/hpcsc/book-stocker-${COMPONENT_NAME}:${VERSION}"

docker buildx build \
  -t "${IMAGE}" \
  --build-arg VERSION="${VERSION}" \
  -f ./Dockerfile.${COMPONENT_NAME} \
  ${BUILD_ARGS} \
  --load \
  .

if [ "${PUSH:-}" = "true" ]; then
  docker push "${IMAGE}"
fi
