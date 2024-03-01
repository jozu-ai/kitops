FROM docker.io/library/golang:alpine AS builder

RUN apk --no-cache upgrade && apk add --no-cache git
ARG version=next
ARG gitCommit=<unknown>

WORKDIR /build

# Cache dependencies in separate layer
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY . .
RUN \
    CGO_ENABLED=0 go build \
    -o kit \
    -ldflags="-s -w -X kitops/pkg/cmd/version.Version=${version} -X kitops/pkg/cmd/version.GitCommit=$gitCommit -X kitops/pkg/cmd/version.BuildTime=$(date -u +'%Y-%m-%dT%H:%M:%SZ')"

FROM docker.io/library/alpine
ENV USER_ID=1001 \
    USER_NAME=kit \
    HOME=/home/user/

RUN apk --no-cache upgrade && \
    apk add --no-cache bash && \
    mkdir -p /home/user/ && \
    adduser -D $USER_NAME -h $HOME -u $USER_ID

USER ${USER_ID}

COPY --from=builder /build/kit /usr/local/bin/kit
COPY build/dockerfiles/KServe/entrypoint.sh /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
