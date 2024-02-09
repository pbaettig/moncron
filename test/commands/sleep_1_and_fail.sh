#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
go run  $SCRIPT_DIR/../../cmd/mon/main.go -name 'do-something-and-fail' -- sh -c 'echo something; sleep 1; exit 127'
