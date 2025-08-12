FROM alpine:latest

# Install system dependencies
RUN apk update && \
    apk add --no-cache curl bash kubectl k9s

# Install clusterctl (latest stable, auto-detect architecture)
ENV CLUSTERCTL_VERSION=v1.10.4
RUN ARCH=$(uname -m) && \
    case ${ARCH} in \
        x86_64) ARCH="amd64" ;; \
        aarch64) ARCH="arm64" ;; \
        armv7l) ARCH="arm" ;; \
        *) echo "Unsupported architecture: ${ARCH}" && exit 1 ;; \
    esac && \
    curl -L https://github.com/kubernetes-sigs/cluster-api/releases/download/${CLUSTERCTL_VERSION}/clusterctl-linux-${ARCH} -o /usr/local/bin/clusterctl && \
    chmod +x /usr/local/bin/clusterctl

# Optional: verify versions
RUN clusterctl version && kubectl version --client && k9s version

ENTRYPOINT [ "/bin/bash" ]
