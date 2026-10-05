FROM docker.io/library/python:3.14-slim

# Mettre à jour les paquets système (corrige ex. CVE-2026-103111 dans libpcre2-8-0)
# et les paquets Python préinstallés vulnérables de l'image de base :
# setuptools (CVE-2025-47273), urllib3 (CVE-2026-97687, CVE-2026-97689),
# msgpack (GHSA-6v7p-g79w-8964). Voir l'onglet Security du repo.
# hadolint ignore=DL3005,DL3042,DL3013
RUN apt-get update && apt-get upgrade -y && \
    apt-get clean && rm -rf /var/lib/apt/lists/* && \
    pip install --no-cache-dir --upgrade pip setuptools urllib3 msgpack

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
