#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd $SCRIPT_DIR/..

dbPath="$SCRIPT_DIR/testdata_1000.db"
# trap "rm -f $dbPath" 0 1 2

if [ ! -f "$dbPath" ]; then
    go run cmd/testdata/*.go -db $dbPath -runs 1000 -hosts 50
fi
echo "DB: $dbPath"
go run cmd/server/*.go -db $dbPath -port 8088 