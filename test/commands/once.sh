#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

name='once'


go run  $SCRIPT_DIR/../../cmd/mon/*.go -name "$name" -server http://localhost:8088/api/runs -stdout -once -- bash -c 'echo "this is an important process that can only run once"; sleep 5.34; ping -c 1 localhost >/dev/null'

