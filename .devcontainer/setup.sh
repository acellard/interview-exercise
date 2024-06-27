#!/bin/bash

# Add docker socker rw permission
sudo chmod o+rw /var/run/docker.sock

# Remplacer le fichier historique de bash par un lien symbolique vers un autre fichier
ln -sfn ${PWD}/.devcontainer/bash_history ~/.bash_history
cat <<'EOF' >> ~/.bashrc
export HISTSIZE=50000
export HISTFILESIZE=100000
EOF

DIST_VSCODE=${HOME}/.vscode-server

if [ $CODESPACES ] ; then
    DIST_VSCODE=${HOME}/.vscode-remote
fi

GENIEAI_HISTORY=${DIST_VSCODE}/data/User/globalStorage/genieai.chatgpt-vscode

# Supprimer GENIEAI_HISTORY existant
rm -rf $GENIEAI_HISTORY

# Créer un nouveau sous-dossier dans .devcontainer
mkdir -p ${PWD}/.devcontainer/genieai.chatgpt-vscode

# Créer un lien symbolique vers ce nouveau sous-dossier
ln -sf ${PWD}/.devcontainer/genieai.chatgpt-vscode $GENIEAI_HISTORY

GOPATH=${HOME}/go

# Supprimer GOPATH existant
rm -rf $GOPATH

# Créer un nouveau sous-dossier dans .devcontainer
mkdir -p ${PWD}/.devcontainer/gopath

# Créer un lien symbolique vers ce nouveau sous-dossier
ln -sf ${PWD}/.devcontainer/gopath $GOPATH
