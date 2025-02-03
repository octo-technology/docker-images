# Référence de l'ensemble des versions à utiliser pour l'installation des outils
FROM docker.io/prom/prometheus:v3.1.0 AS prometheus
FROM docker.io/prom/alertmanager:v0.28.0 AS alertmanager
FROM public.ecr.aws/aws-cli/aws-cli:2.23.11 AS aws-cli
FROM docker.io/library/docker:27 AS docker-compose
FROM docker.io/bitnami/kubectl:1.32 AS kubectl
FROM docker.io/bitnami/trivy:0.58.2 AS trivy
FROM docker.io/derailed/k9s:v0.32.7 AS k9s
FROM docker.io/wagoodman/dive:v0.12 AS dive

FROM ghcr.io/coder/code-server:4.96.2-ubuntu
ARG KUBECTX_VERSION="v0.9.5"

# Nécessaire pour installer les paquets
USER root
WORKDIR /home/coder
# Mettre à jour les paquets et installer les dépendances
RUN export DEBIAN_FRONTEND=noninteractive ; export TZ=Europe/Paris ; apt-get -y update && apt-get -y install \
    git vim wget htop amazon-ecr-credential-helper groff-base gnupg2 pass bind9-utils gettext-base apache2-utils && \
    # Nettoyer les fichiers inutiles
    apt-get clean && rm -rf /var/lib/apt/lists/*
# Et on revient à un utilisateur lambda pour la suite
USER coder

# Install AWS cli
COPY --from=aws-cli /usr/local /usr/local

# Installer docker-compose extension
COPY --from=docker-compose /usr/local/bin/docker-compose /usr/local/lib/docker/cli-plugins/docker-compose

# Installer kubectl/trivy/k9s/dive
COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /usr/local/bin
COPY --from=trivy /opt/bitnami/trivy/bin/trivy /usr/local/bin
COPY --from=k9s /bin/k9s /usr/local/bin
COPY --from=dive /usr/local/bin/dive /usr/local/bin

# Installer kubectx
ADD --chmod=755 https://github.com/ahmetb/kubectx/releases/download/${KUBECTX_VERSION}/kubectx /usr/local/bin/kubectx
ADD --chmod=755 https://github.com/ahmetb/kubectx/releases/download/${KUBECTX_VERSION}/kubens /usr/local/bin/kubens

# Installer kube-ps1
COPY bash_completion.sh /etc/bash_completion
ADD --chmod=755 https://raw.githubusercontent.com/jonmosco/kube-ps1/refs/heads/master/kube-ps1.sh /usr/local/bin/kube-ps1.sh

# Installer promtool + prometheus + alertmanager + amtool
COPY --from=prometheus /bin/prometheus /usr/local/bin
COPY --from=prometheus /bin/promtool /usr/local/bin
COPY --from=alertmanager /bin/alertmanager /usr/local/bin
COPY --from=alertmanager /bin/amtool /usr/local/bin

# Définir le point d'entrée par défaut
CMD ["/bin/bash"]
