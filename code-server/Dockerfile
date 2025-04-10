# Référence de l'ensemble des versions à utiliser pour l'installation des outils
FROM docker.io/prom/prometheus:v3.2.1 AS prometheus
FROM docker.io/prom/alertmanager:v0.28.1 AS alertmanager
FROM public.ecr.aws/aws-cli/aws-cli:2.26.0 AS aws-cli
FROM docker.io/library/docker:28.0.4-cli AS docker
FROM docker.io/docker/buildx-bin:0.22.0 AS buildx
FROM docker.io/bitnami/kubectl:1.32 AS kubectl
FROM docker.io/bitnami/trivy:0.61.0 AS trivy
FROM docker.io/derailed/k9s:v0.50.2 AS k9s
FROM ghcr.io/wagoodman/dive:v0.13.1 AS dive
FROM ghcr.io/grafana/k6:0.58.0 AS k6
FROM docker.io/alpine/helm:3.17.2 AS helm
FROM docker.io/hadolint/hadolint:v2.12.0-alpine AS hadolint
FROM docker.io/goodwithtech/dockle:v0.4.15 AS dockle
FROM gcr.io/k8s-skaffold/skaffold:v2.15.0 AS skaffold
# A opentofu minimal image will be published soon : https://github.com/opentofu/opentofu/issues/2356
FROM ghcr.io/opentofu/opentofu:1.9.0 AS opentofu
FROM docker.io/hashicorp/terraform:1.11.4 AS terraform
FROM registry.gitlab.com/gitlab-org/cli:v1.55.0 AS gitlab-cli

FROM ghcr.io/coder/code-server:4.99.1-noble
ARG TARGETOS
ARG TARGETARCH

# Nécessaire pour installer les paquets
USER root
WORKDIR /home/coder
ENV DEBIAN_FRONTEND=noninteractive TZ=Europe/Paris
# Mettre à jour les paquets et installer les dépendances
RUN apt-get -y update && apt-get -y install \
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
    # Nettoyer les fichiers inutiles
    && apt-get clean && rm -rf /var/lib/apt/lists/* && \
    rm -f /etc/ssl/private/ssl-cert-snakeoil.key /etc/ssl/certs/ssl-cert-snakeoil.pem
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
COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin
COPY --from=trivy /opt/bitnami/trivy/bin/trivy /usr/local/bin
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
ENV KREW_VERSION="v0.4.4"
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
ENV KAPP_VERSION="v0.64.0"
ADD --chmod=755 --chown=root:root https://github.com/carvel-dev/kapp/releases/download/${KAPP_VERSION}/kapp-${TARGETOS}-${TARGETARCH} /usr/local/bin/kapp

# Installer kubectx
ENV KUBECTX_VERSION="v0.9.5"
ADD --chmod=755 --chown=root:root https://github.com/ahmetb/kubectx/releases/download/${KUBECTX_VERSION}/kubectx /usr/local/bin/kubectx
ADD --chmod=755 --chown=root:root https://github.com/ahmetb/kubectx/releases/download/${KUBECTX_VERSION}/kubens /usr/local/bin/kubens

# Installer kube-ps1
COPY bash_completion.sh /etc/bash_completion.d/kube.sh
ADD --chmod=755 https://raw.githubusercontent.com/jonmosco/kube-ps1/refs/heads/master/kube-ps1.sh /usr/local/bin/kube-ps1.sh

# Installer promtool + prometheus + alertmanager + amtool
COPY --from=prometheus /bin/prometheus /usr/local/bin
COPY --from=prometheus /bin/promtool /usr/local/bin
COPY --from=alertmanager /bin/alertmanager /usr/local/bin
COPY --from=alertmanager /bin/amtool /usr/local/bin

# Installation de plugins ❤️
RUN code-server --install-extension bierner.markdown-mermaid

# Définir le point d'entrée par défaut
CMD ["/bin/bash"]
