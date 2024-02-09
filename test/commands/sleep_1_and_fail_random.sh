#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
go run  $SCRIPT_DIR/../../cmd/mon/main.go -name 'do-something' -- bash -c 'echo something; sleep 1; exit $(( RANDOM % 127 ))'
