#!/bin/bash

SCRIPT_DIR=$(dirname "${BASH_SOURCE[0]}")
PROJECT_DIR=$(realpath "$SCRIPT_DIR/..")
CMD_DIR="${PROJECT_DIR}/cmd"

cd ${PROJECT_DIR}
for CMD in `ls $CMD_DIR`; do
  go build ./cmd/$CMD
done
