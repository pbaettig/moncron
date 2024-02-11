#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
# go run  $SCRIPT_DIR/../../cmd/mon/*.go \
./mon \
    -name 'random-rss' \
    -pushgw 'http://localhost:9091' \
    -web 'http://localhost:8080' \
    -log "/tmp/output.log" \
    -log-size $((12*1024)) \
    -timeout 1s \
    -- \
    python3 -c 'import random; import sys; print(sum([n for n in range(random.randint(10_000_000, 50_000_000))])); sys.exit(random.randint(0,1))'
echo $?