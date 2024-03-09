#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
go run  $SCRIPT_DIR/../../cmd/mon/*.go -name 'sleep-random' -stdout -- python3 -c 'import time; import random; time.sleep((random.random() * 4)+1)'
