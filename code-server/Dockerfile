ARG CODE_SERVER_VERSION=4.93.1
ARG UBUNTU_CODENAME=noble

# We use the official code-server image as base, in Ubuntu 24.04 (noble)
FROM ghcr.io/coder/code-server:$CODE_SERVER_VERSION-${UBUNTU_CODENAME}

# Install the required packages
USER root
RUN apt-get update && apt-get install -y \
    make=4.3-4.1build2 \
    procps=2:4.0.4-4ubuntu3.1 \
    curl=8.5.0-2ubuntu10.4 \
    file=1:5.45-3build1 \
    git=1:2.43.0-1ubuntu7.1 \
    vim=2:9.1.0016-1ubuntu7.2 \
    wget=1.21.4-1ubuntu4.1 \
    htop=3.3.0-4build1 \
    dnsutils \
    gnupg2 \
    pass \
    j2cli \
    gettext-base \
    dnsutils=1:9.18.28-0ubuntu0.24.04.1 \
    iputils-ping=3:20240117-1build1 \
    unzip=6.0-28ubuntu4 \
    zip=3.0-13build1 \
    postgresql-client \
    redis-tools \
    mariadb-client \
    && rm -rf /var/lib/apt/lists/*

# Instal AWSCLI
RUN curl -sL https://awscli.amazonaws.com/awscli-exe-linux-$(lscpu | awk '/Architecture:/{print $2}').zip -o awscliv2.zip \
    && unzip awscliv2.zip \
    && ./aws/install \
    && rm -rf awscliv2.zip ./aws

# Install docker cli & docker-compose
ADD --chmod=644 https://download.docker.com/linux/ubuntu/gpg /etc/apt/keyrings/docker.asc
RUN echo \
        "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
        $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
        sudo tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt-get update && apt-get install -y \
        docker-ce-cli \
        docker-compose-plugin \
    && rm -rf /var/lib/apt/lists/*

# Install amazon-ecr-credential-helper
# This package is required to authenticate with AWS ECR registries and must be 
# installed AFTER the docker cli to avoid conflicts between docker.io and docker-ce-cli packages
RUN apt-get update && apt-get install -y \
    amazon-ecr-credential-helper=0.7.1-1ubuntu0.24.04.1 \
    && rm -rf /var/lib/apt/lists/*

# Install dive (needed for DKUSR & DKADM training)
ENV DIVE_VERSION=0.12.0
RUN curl -sL https://github.com/wagoodman/dive/releases/download/v${DIVE_VERSION}/dive_${DIVE_VERSION}_linux_$(dpkg --print-architecture).deb -o dive.deb \
    && apt install ./dive.deb \
    && rm dive.deb

# Install kubectl (needed for DKUSR & DKADM training)
ENV KUBECTL_VERSION=1.30.3
RUN wget https://dl.k8s.io/release/v${KUBECTL_VERSION}/bin/linux/$(dpkg --print-architecture)/kubectl \
    && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl \
    && rm kubectl \
    && kubectl completion bash > /etc/bash_completion.d/kubectl \
    && echo 'alias k=kubectl' >> /root/.bashrc

# Install helm (needed for DKUSR & DKADM training)
ENV HELM_VERSION=3.14.4
RUN wget https://get.helm.sh/helm-v${HELM_VERSION}-linux-$(dpkg --print-architecture).tar.gz \
    && tar -zxvf helm-v${HELM_VERSION}-linux-$(dpkg --print-architecture).tar.gz \
    && install -o root -g root -m 0755 linux-$(dpkg --print-architecture)/helm /usr/local/bin/helm \
    && rm -rf helm-v${HELM_VERSION}-linux-$(dpkg --print-architecture).tar.gz linux-$(dpkg --print-architecture)

# Install k9s (needed for DKUSR & DKADM training)
ENV K9S_VERSION=0.32.4
RUN curl -sL https://github.com/derailed/k9s/releases/download/v${K9S_VERSION}/k9s_linux_$(dpkg --print-architecture).deb -o k9s.deb \
    && dpkg -i k9s.deb \
    && rm k9s.deb

# Install krew (needed for DKUSR & DKADM training)
ENV KREW_VERSION=0.4.1
RUN curl -sL "https://github.com/kubernetes-sigs/krew/releases/download/v${KREW_VERSION}/krew-linux_$(dpkg --print-architecture).tar.gz" | tar zxvf - \
    && ./krew-"$(dpkg --print-architecture)" install krew \
    && rm -rf krew.tar.gz krew-"$(dpkg --print-architecture)"
ENV PATH="${PATH}:/root/.krew/bin"

# Install krew plugins : kubectx, kubens & neat
RUN kubectl krew install ctx ns neat

# Install kube-ps1 (needed for DKUSR & DKADM training)
ENV KUBE_PS1_VERSION=0.7.0
RUN git clone --depth=1 git@github.com:jonmosco/kube-ps1.git /opt/kube-ps1 \
    && echo 'source /etc/profile.d/kube-ps1.sh' >> /root/.bashrc \
    && echo "export PS1='[${debian_chroot:+($debian_chroot)}\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\] $(kube_ps1)] \$ '" >> /root/.bashrc

USER ubuntu
