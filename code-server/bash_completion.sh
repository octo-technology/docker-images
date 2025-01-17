source /usr/local/bin/kube-ps1.sh
export PS1='[${debian_chroot:+($debian_chroot)}\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\] $(kube_ps1)] $ '
alias aws_ecr_login='aws ecr get-login-password --region $AWS_REGION | docker login --username $DOCKER_REGISTRY_USERNAME --password-stdin $REGISTRY_URL'
