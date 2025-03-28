#!/bin/bash

source <(docker completion bash)
source <(kubectl completion bash)
alias k=kubectl
complete -o default -F __start_kubectl k

source /usr/local/bin/kube-ps1.sh
export KUBE_PS1_PREFIX=" ("
export KUBE_PS1_ENABLED=off
export PS1='[\[\033[01;32m\]\u@training\[\033[00m\]:\[\033[01;34m\]\w\[\033[00m\]$(kube_ps1)] $ '
