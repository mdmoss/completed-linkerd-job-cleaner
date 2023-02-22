#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

VERSION=$1

docker buildx build .\
  --platform linux/amd64 \
  --tag mdmoss/completed-linkerd-job-cleaner:$VERSION \
  --push