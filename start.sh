#!/bin/bash
kill -9 `pgrep tg-bot-linux`
BIN_DIR=$(cd `dirname $0`; pwd)
$BIN_DIR/tg-bot-linux $BIN_DIR/config.json > $BIN_DIR/run-`date +%s`.log 2>&1 &

