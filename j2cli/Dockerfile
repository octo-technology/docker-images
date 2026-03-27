FROM docker.io/library/python:3.14-slim

WORKDIR /app

# Labels OCI pour identification de l'image
ARG VERSION="dev"
ARG BUILD_TIMESTAMP="unknown"

LABEL org.opencontainers.image.title="J2CLI"
LABEL org.opencontainers.image.authors="Joseph Page <josephpage@users.noreply.github.com>, Yannig Perre <yannig@users.noreply.github.com>"
LABEL org.opencontainers.image.description="Jinja2 CLI tool for rendering templates from the command line"
LABEL org.opencontainers.image.url="https://www.octo.academy"
LABEL org.opencontainers.image.vendor="OCTO Technology"
LABEL org.opencontainers.image.source="https://github.com/octo-technology/docker-images/tree/main/j2cli"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.created="${BUILD_TIMESTAMP}"
LABEL org.opencontainers.image.base.name="docker.io/library/python:3.13-slim"

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

ENTRYPOINT ["/usr/bin/bash"]
