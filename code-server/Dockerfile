# Référence de l'ensemble des versions à utiliser pour l'installation des outils
FROM docker.io/prom/prometheus:v3.9.1 AS prometheus
FROM docker.io/prom/alertmanager:v0.31.1 AS alertmanager
FROM public.ecr.aws/aws-cli/aws-cli:2.33.21 AS aws-cli
FROM docker.io/library/docker:29.2.1-cli AS docker
FROM docker.io/docker/buildx-bin:0.31.1 AS buildx
FROM docker.io/rancher/kubectl:v1.35.0 AS kubectl
FROM ghcr.io/aquasecurity/trivy:0.69.1 AS trivy
FROM docker.io/derailed/k9s:v0.50.18 AS k9s
FROM ghcr.io/wagoodman/dive:v0.13.1 AS dive
FROM ghcr.io/grafana/k6:1.6.0 AS k6
FROM docker.io/alpine/helm:4.1.1 AS helm
FROM docker.io/hadolint/hadolint:v2.14.0-alpine AS hadolint
FROM docker.io/goodwithtech/dockle:v0.4.15 AS dockle
FROM gcr.io/k8s-skaffold/skaffold:v2.17.2 AS skaffold
FROM ghcr.io/opentofu/opentofu:1.11.5-minimal AS opentofu
FROM docker.io/hashicorp/terraform:1.14.5 AS terraform
FROM registry.gitlab.com/gitlab-org/cli:v1.85.0 AS gitlab-cli

FROM ghcr.io/coder/code-server:4.108.2-noble
ARG TARGETOS
ARG TARGETARCH

# Nécessaire pour installer les paquets
USER root
WORKDIR /home/coder
ENV DEBIAN_FRONTEND=noninteractive TZ=Europe/Paris
# Mettre à jour les paquets et installer les dépendances
RUN apt-get -y update && apt-get -y install --no-install-recommends \
    git \
    vim \
    wget \
    htop \
    unzip \
    bash-completion \
    groff-base \
    gnupg2 \
    pass \
    bind9-utils \
    gettext-base \
    apache2-utils \
    emacs \
    jq \
    ipcalc \
    xz-utils \
    gron \
    python3-pip \
    python3-venv \
    # Nettoyer les fichiers inutiles
    && apt-get clean && rm -rf /var/lib/apt/lists/* && \
    rm -f /etc/ssl/private/ssl-cert-snakeoil.key /etc/ssl/certs/ssl-cert-snakeoil.pem && \
    # Création d'alias python/pip => python3/pip3
    ln -s /usr/bin/python3 /usr/local/bin/python && \
    ln -s /usr/bin/pip3 /usr/local/bin/pip

# Et on revient à un utilisateur lambda pour la suite
USER coder

# Install AWS CLI & docker-credential-ecr-login
COPY --from=aws-cli /usr/local /usr/local
ENV AWS_ECR_HELPER_VERSION="0.9.1"
ADD --chmod=755 --chown=root:root \
    https://amazon-ecr-credential-helper-releases.s3.us-east-2.amazonaws.com/${AWS_ECR_HELPER_VERSION}/${TARGETOS}-${TARGETARCH}/docker-credential-ecr-login \
    /usr/local/bin/docker-credential-ecr-login

# Installer docker & docker-compose
COPY --from=docker /usr/local/bin/docker /usr/local/bin/docker
COPY --from=docker /usr/local/bin/docker-compose /usr/local/lib/docker/cli-plugins/docker-compose
COPY --from=buildx /buildx /usr/local/lib/docker/cli-plugins/docker-buildx

# Installer cli tools
COPY --from=kubectl /bin/kubectl /usr/local/bin
COPY --from=trivy /usr/local/bin/trivy /usr/local/bin
COPY --from=k9s /bin/k9s /usr/local/bin
COPY --from=dive /usr/local/bin/dive /usr/local/bin
COPY --from=k6 /usr/bin/k6 /usr/local/bin
COPY --from=helm /usr/bin/helm /usr/local/bin
COPY --from=hadolint /bin/hadolint /usr/local/bin
COPY --from=dockle /usr/bin/dockle /usr/local/bin
COPY --from=skaffold /usr/bin/skaffold /usr/local/bin
COPY --from=opentofu /usr/local/bin/tofu /usr/local/bin
COPY --from=terraform /bin/terraform /usr/local/bin
COPY --from=gitlab-cli /usr/bin/glab /usr/local/bin

# Install krew
ENV KREW_VERSION="v0.4.5"
# hadolint ignore=DL3003,SC2016
RUN KREW_BINARY="krew-${TARGETOS}_${TARGETARCH}" && \
    TMPDIR="$(mktemp -d)" && \
    cd "$TMPDIR" && \
    wget -q https://github.com/kubernetes-sigs/krew/releases/download/${KREW_VERSION}/${KREW_BINARY}.tar.gz && \
    tar zxvf "${KREW_BINARY}.tar.gz" && \
    ./${KREW_BINARY} install krew && \
    echo 'PATH="$HOME/.krew/bin:$PATH"' >> ~/.profile && \
    cd ~ && rm -rf "$TMPDIR" && rm -f ~/.wget-hsts

# Install kapp
ENV KAPP_VERSION="v0.64.2"
ADD --chmod=755 --chown=root:root https://github.com/carvel-dev/kapp/releases/download/${KAPP_VERSION}/kapp-${TARGETOS}-${TARGETARCH} /usr/local/bin/kapp

# Installer kubectx
ENV KUBECTX_VERSION="v0.9.5"
ADD --chmod=755 --chown=root:root https://github.com/ahmetb/kubectx/releases/download/${KUBECTX_VERSION}/kubectx /usr/local/bin/kubectx
ADD --chmod=755 --chown=root:root https://github.com/ahmetb/kubectx/releases/download/${KUBECTX_VERSION}/kubens /usr/local/bin/kubens

# Installer kube-ps1
COPY bash_completion.sh /etc/bash_completion.d/kube.sh
ADD --chmod=755 https://raw.githubusercontent.com/jonmosco/kube-ps1/refs/heads/master/kube-ps1.sh /usr/local/bin/kube-ps1.sh

# Installer promtool + prometheus + alertmanager + amtool
COPY --from=prometheus /bin/prometheus /bin/promtool /usr/local/bin/
COPY --from=alertmanager /bin/alertmanager /bin/amtool /usr/local/bin/

# Installation de plugins ❤️
RUN code-server --install-extension bierner.markdown-mermaid

# Définir le point d'entrée par défaut
CMD ["/bin/bash"]
