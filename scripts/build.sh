#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR/..

build_time="$(date '+%d.%m.%y %T')"
version="latest"
if [[ ! -z "$CI" ]]; then
    version="$GITHUB_REF_NAME-${GITHUB_SHA:0:7}"
fi

CGO_ENABLED=0 go build \
    -o mon \
    --ldflags="-X 'github.com/pbaettig/moncron/internal/pkg/buildinfo.Version=$version'" \
    cmd/mon/*.go

# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
#     -o mon \
#     --ldflags="-X 'github.com/pbaettig/moncron/internal/buildinfo.Version=$version'" \
#     cmd/mon/*.go