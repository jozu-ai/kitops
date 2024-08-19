# Containerized Kit CLI

The Dockerfile in this directory can be used for building a containerized version of the Kit CLI. To build this image manually, run the following command from the root of this repository:
```bash
docker build -t kit-cli:latest -f build/dockerfiles/Dockerfile .
```

The default entrypoint for this image is invoking the `kit` CLI, allowing for quickly running Kit commands without needing to install the CLI:
```
docker run --rm ghcr.io/jozu-ai/kit:latest version
```
