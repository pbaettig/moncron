#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
go run  $SCRIPT_DIR/../../cmd/mon/main.go -name 'sleep-random' -- python3 -c 'import time; import random; time.sleep((random.random() * 4)+1)'
