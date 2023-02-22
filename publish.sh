#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

VERSION=$1
MAJOR_VERSION=$(echo $VERSION | sed -r 's/([0-9]+)\..*/\1/')

read -p "Publish version $VERSION (plus major version $MAJOR_VERSION) to Docker Hub? " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    docker buildx build .\
        --platform linux/amd64 \
        --tag mdmoss/completed-linkerd-job-cleaner:$VERSION \
        --tag mdmoss/completed-linkerd-job-cleaner:$MAJOR_VERSION \
        --push
fi
