#!/bin/bash

set -e

# Kit configuration via env vars:
#   - UNPACK_PATH   : Where to unpack the modelkit (default: /home/user/modelkit)
#   - UNPACK_FILTER : What to unpack from the modelkit (default: model)
#   - MODELKIT_REF  : The modelkit to unpack -- required!
UNPACK_PATH=${UNPACK_PATH:-/home/user/modelkit/}
UNPACK_FILTER=${UNPACK_FILTER:-model}
if [ -z "$MODELKIT_REF" ]; then
  echo "Environment variable \$MODELKIT_REF is required"
  exit 1
fi

# Variables for verifying signature via cosign. Can specify either a key to use for
# verifying or an identity and oidc issuer for keyless verification
if [ -n "$COSIGN_KEY" ]; then
  echo "Verifying signature for modelkit $MODELKIT_REF via key"
  cosign verify --key "$COSIGN_KEY" "$MODELKIT_REF"
elif [ -n "$COSIGN_CERT_IDENTITY" ] && [ -n "$COSIGN_CERT_OIDC_ISSUER" ]; then
  echo "Verifying signature for modelkit $MODELKIT_REF"
  cosign verify \
    --certificate-identity=${COSIGN_CERT_IDENTITY} \
    --certificate-oidc-issuer=${COSIGN_CERT_OIDC_ISSUER} \
    "$MODELKIT_REF"
else
  echo "Signature verification is disabled"
fi

echo "Binary version info:"
kit version

echo "Unpacking modelkit $MODELKIT_REF to $UNPACK_PATH with filter $UNPACK_FILTER"
kit unpack "$MODELKIT_REF" --dir "$UNPACK_PATH" --filter="$UNPACK_FILTER"
