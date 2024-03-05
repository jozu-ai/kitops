#!/bin/bash

set -e

# Find directory this script is in and work relative to that to make the script
# callable from anywhere
SCRIPT_DIR=$(cd "$(dirname "$0")" || exit; pwd)
DOCS_DIR="$SCRIPT_DIR"

echo "Generating base files..."
pushd "$DOCS_DIR" >/dev/null
go run gen-doc.go
popd >/dev/null

echo "Generated files:"
find "$DOCS_DIR/" -name "kit*.md"
echo ""

# Truncate cli-reference doc and insert header
cat "$DOCS_DIR/cli-reference.header" > "$DOCS_DIR/cli-reference.md"

echo "Building cli-reference.md"
for file in $(find "$DOCS_DIR/" -name "kit_*.md" | sort); do
  # Trim off all "See also" sections from each command before adding to doc
  sed -n '/### SEE ALSO/q;p' "$file" >> $DOCS_DIR/cli-reference.md
done

# Remove generated files, keeping only the combined CLI reference doc
rm -rf $DOCS_DIR/kit*.md
