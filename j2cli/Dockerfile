FROM docker.io/library/python:3.14-slim

# hadolint ignore=DL3013
RUN apt-get update && apt-get upgrade -y && \
    apt-get clean && rm -rf /var/lib/apt/lists/* && \
    pip install --no-cache-dir --upgrade pip setuptools urllib3 msgpack

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
