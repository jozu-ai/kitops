#!/bin/bash

set -e

# R2 Configuration from environment variables
R2_ACCESS_KEY_ID="${R2_ACCESS_KEY_ID}"
R2_SECRET_ACCESS_KEY="${R2_SECRET_ACCESS_KEY}"
R2_ENDPOINT="${R2_ENDPOINT}"

R2_BUCKET_NAME="kitops-binaries"
REMOTE_NAME="r2"

# Files to Upload
FILES_TO_UPLOAD=(
    "llamafile.tar.gz"
    "ui.tar.gz"
    "checksums.txt"
)


SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
LLAMAFILE_VER=$(go run $SCRIPT_DIR/llamafile_ver_helper.go)



# ---------------------------
# Functions
# ---------------------------

configure_rclone() {
    echo "Configuring rclone for Cloudflare R2..."

    # Create rclone configuration using environment variables
    rclone config create "$REMOTE_NAME" s3 \
        provider Cloudflare \
        access_key_id "$R2_ACCESS_KEY_ID" \
        secret_access_key "$R2_SECRET_ACCESS_KEY" \
        endpoint "$R2_ENDPOINT" \
        region auto \
        acl private

    echo "rclone configured successfully."
}

upload() {
    SOURCE="$1"
    DESTINATION="$2"

    FULL_DESTINATION="${REMOTE_NAME}:${R2_BUCKET_NAME}/${DESTINATION}"
    
    # Perform the upload using rclone copy
    echo "Uploading '${SOURCE}' to '${FULL_DESTINATION}' ..."
    rclone copy "$SOURCE" "$FULL_DESTINATION" --verbose
    echo "Upload of '${SOURCE}' completed successfully."
}

# ---------------------------
# Script Execution
# ---------------------------

# Check if rclone is installed
if ! command -v rclone &> /dev/null
then
    echo "Error: rclone could not be found. Please install rclone before running this script."
    exit 1
fi

# Check if all required environment variables are set
if [[ -z "$R2_ACCESS_KEY_ID" || -z "$R2_SECRET_ACCESS_KEY" || -z "$R2_ENDPOINT" ]]; then
    echo "Error: One or more required environment variables (R2_ACCESS_KEY_ID, R2_SECRET_ACCESS_KEY, R2_ENDPOINT) are not set."
    exit 1
fi

configure_rclone

for DEST_PATH in "${FILES_TO_UPLOAD[@]}"; do
    FILE_PATH="${SCRIPT_DIR}/../${DEST_PATH}"
    
    if [ -f "$FILE_PATH" ]; then
        echo "File found: '$FILE_PATH'. Preparing to upload..."
        upload "$FILE_PATH" "$LLAMAFILE_VER/$DEST_PATH"
    else
        echo "Warning: File '$FILE_PATH' does not exist. Skipping upload."
    fi
done

echo "All files uploaded"
