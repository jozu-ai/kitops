#!/bin/bash

set -e

echo "Binary version info:"
kit version

read -r -a UNPACK_FLAGS <<< "$KIT_UNPACK_FLAGS"

if [ $# != 2 ]; then
  echo "Usage: entrypoint.sh <src-uri> <dest-path>"
  exit 1
fi

REPO_NAME="${1#kit://}"
OUTPUT_DIR="$2"

echo "Unpacking $REPO_NAME to $OUTPUT_DIR"
echo "Unpack options: ${KIT_UNPACK_FLAGS}"
kit unpack "$REPO_NAME" -d "$OUTPUT_DIR" "${UNPACK_FLAGS[@]}"

echo "Unpacked modelkit:"
cat "$OUTPUT_DIR/Kitfile"
