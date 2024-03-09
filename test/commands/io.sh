#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

outfile="$(mktemp)"
trap "rm -f $outfile" 0 1 2

go run  $SCRIPT_DIR/../../cmd/mon/*.go \
    -name 'io' \
    -pushgw 'http://localhost:9091' \
    -server 'http://localhost:8088/api/runs' \
    -stdout \
    -log "/tmp/output.log" \
    -log-size $((12*1024)) \
    -- \
    dd if=/dev/urandom of=$outfile bs=4096 count=120000
echo $?