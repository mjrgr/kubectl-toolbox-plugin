#!/bin/sh
set -u
set -e
export MSYS_NO_PATHCONV=1
current_path=`( cd \`dirname "$0"\` && pwd )`

podman build . -t kubectl-toolbox-plugin
podman container run -it --rm -v "${HOME}/.kube:/home/dummy/.kube" kubectl-toolbox-plugin $@