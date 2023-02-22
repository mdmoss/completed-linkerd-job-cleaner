#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

VERSION=$1

docker buildx build .\
  --platform linux/amd64 \
  --tag mdmoss/linkerd-completed-job-cleaner:$VERSION \
  --push