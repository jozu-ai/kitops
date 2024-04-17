# Functional testing structure

This directory contains subdirectories for various test cases for the Kit CLI. Each subdirectory is named based on what flow is being tested, and contains files describing a test case via the struct

```yaml
description: "Description of what the test is testing"
kitfile: |
  The Kitfile to use for the test
kitignore: |
  The .kitignore to use for the test
files:
  - A list of files that should be included in the modelkit
  - E.g. "dir1/dir2/myfile.txt"
  - Directories will be created as necessary
ignored:
  - A list of files that should exist in the context but should _not_ be included in the modelkit
```

Fields are optional and their usage depends on the test -- e.g. testing packing and unpacking will require `files:` and `ignored:`, whereas testing tagging does not require files necessarily.

To generate a test case from an existing modelkit directory, you can use the following snippet:
```bash
cat <<EOF > new-test-case.yaml
description: "enter your description"
kitfile: |
$(sed 's|^|  |g' Kitfile)
kitignore: |
$(sed 's|^|  |g' .kitignore)
files: # Sort these into files vs ignored as necessary
$(find . -type f | sed 's|^|  - |g')
EOF
```
