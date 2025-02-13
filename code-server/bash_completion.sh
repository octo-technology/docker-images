#!/bin/bash

source /usr/share/bash-completion/completions/git
source <(docker completion bash)
source <(kubectl completion bash)
alias k=kubectl
complete -o default -F __start_kubectl k

source /usr/local/bin/kube-ps1.sh
export PS1='[${debian_chroot:+($debian_chroot)}\[\033[01;32m\]\u@\h\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\] $(kube_ps1)] $ '
