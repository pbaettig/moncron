#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
go run  $SCRIPT_DIR/../../cmd/mon/main.go -verbose -name 'random-rss' -- python3 -c 'import random; print(sum([n for n in range(random.randint(1_000_000, 50_000_000))]))'
