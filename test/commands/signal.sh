#!/bin/bash
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

name='wait-for-signal'

go run  $SCRIPT_DIR/../../cmd/mon/*.go -name "$name" -- sleep 42.12 &

sleep 1

ps -ef | grep 'sleep 42.12' | grep -v grep | grep -v 'go' | grep -v "$name" | tail -1 | awk '{print $2}' | xargs kill

