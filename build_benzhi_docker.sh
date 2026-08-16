#!/usr/bin/env bash
set -euo pipefail

IMAGE_NAME="${1:-pet-benzhi}"
IMAGE_TAG="${2:-latest}"

docker build \
  -f benzhi.Dockerfile \
  -t "${IMAGE_NAME}:${IMAGE_TAG}" \
  .
