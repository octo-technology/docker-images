FROM docker.io/bats/bats:1.13.0

# Labels OCI pour identification de l'image
ARG VERSION="dev"
ARG BUILD_TIMESTAMP="unknown"

LABEL org.opencontainers.image.title="BATS"
LABEL org.opencontainers.image.authors="Joseph Page <josephpage@users.noreply.github.com>, Yannig Perre <yannig@users.noreply.github.com>"
LABEL org.opencontainers.image.description="BATS (Bash Automated Testing System) with curl, used in OCTO training sessions"
LABEL org.opencontainers.image.url="https://www.octo.academy"
LABEL org.opencontainers.image.vendor="OCTO Technology"
LABEL org.opencontainers.image.source="https://github.com/octo-technology/docker-images/tree/main/bats"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.created="${BUILD_TIMESTAMP}"
LABEL org.opencontainers.image.base.name="docker.io/bats/bats:1.13.0"

RUN apk update && apk upgrade --no-cache && \
    apk add --no-cache curl
