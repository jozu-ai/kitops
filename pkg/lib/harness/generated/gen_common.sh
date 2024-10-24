#!/bin/sh

UI_DIR=../../../../frontend/dev-mode
LLAMAFILE_VER=$(go run ./llamafile_ver_helper.go)

build_ui() {
    UI_HOME=$1

    echo "Building harness UI from ${UI_HOME}"
    
    pushd ${UI_HOME}
    pnpm install
    pnpm run build
    popd
}

compress() {
    echo "Compressing $1 to $2"
    tar -czf $2 -C $1 .
}

generate_sha() {
    FILE=$1
    CHECKSUM_FILE=$2

    echo "Generating SHA256 checksum for ${FILE}"
    sha256sum ${FILE} | awk '{print $1 "  " FILENAME}' FILENAME="${FILE}" >> ${CHECKSUM_FILE}
    echo "Checksum for ${FILE} saved to ${CHECKSUM_FILE}"
}


# Function to download a binary asset from a GitHub release
download_github_release_binary() {
    OWNER=$1
    REPO=$2
    RELEASE_TAG=$3
    ASSET_NAME=$4
    OUTPUT_DIR=$5

    echo "Downloading asset '${ASSET_NAME}' from release '${RELEASE_TAG}' of repository '${OWNER}/${REPO}'"

    DOWNLOAD_URL="https://github.com/${OWNER}/${REPO}/releases/download/${RELEASE_TAG}/llamafile-${RELEASE_TAG}"
    
    mkdir -p ${OUTPUT_DIR}

    curl -L -o ${OUTPUT_DIR}/llamafile ${DOWNLOAD_URL}

    echo "Asset '${ASSET_NAME}' downloaded successfully to: ${OUTPUT_DIR}/${ASSET_NAME}"

    echo "${RELEASE_TAG}" > ${OUTPUT_DIR}/llamafile.version

    COMPRESSED_FILE="llamafile.tar.gz"
    compress ${OUTPUT_DIR} ../${COMPRESSED_FILE}

    echo "Compressed asset saved to: ${COMPRESSED_FILE}"
}

build_ui ${UI_DIR}
compress ${UI_DIR}/dist ../ui.tar.gz
download_github_release_binary "Mozilla-Ocho" "llamafile" "${LLAMAFILE_VER}" "llamafile-${LLAMAFILE_VER}" "./downloads"

CHECKSUM_FILE="../checksums.txt"
> ${CHECKSUM_FILE}  # Clear the checksum file if it exists
generate_sha "../ui.tar.gz" ${CHECKSUM_FILE}
generate_sha "../llamafile.tar.gz" ${CHECKSUM_FILE}

echo "All checksums have been saved to ${CHECKSUM_FILE}"