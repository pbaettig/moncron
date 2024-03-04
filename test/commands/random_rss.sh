#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
go run  $SCRIPT_DIR/../../cmd/mon/*.go \
    -name 'random-rss' \
    -pushgw 'http://localhost:9091' \
    -server 'http://localhost:8088/api/runs' \
    -log "/tmp/output.log" \
    -log-size $((12*1024)) \
    -timeout 1s \
    -- \
    python3 -c 'import random; import sys; print(sum([n for n in range(random.randint(1_000_000, 2_000_000))])); sys.exit(random.randint(0,1))'
echo $?