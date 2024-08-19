#!/bin/bash

set -e

UNPACK_PATH=${UNPACK_PATH:-/home/user/modelkit/}
UNPACK_FILTER=${UNPACK_FILTER:-model}
if [ -z "$MODELKIT_REF" ]; then
  echo "Environment variable \$MODELKIT_REF is required"
  exit 1
fi

echo "Binary version info:"
kit version

echo "Unpacking modelkit $MODELKIT_REF to $UNPACK_PATH with filter $UNPACK_FILTER"
kit unpack "$MODELKIT_REF" --dir "$UNPACK_PATH" --filter="$UNPACK_FILTER"
