#!/usr/bin/env bash
set -euo pipefail

# ——————————————————————————————————————
# Variables globales
# ——————————————————————————————————————
BIN_DIR="/usr/local/bin"
SERVICE_NAME="justdeploy.service"
SYSTEMD_DIR="/etc/systemd/system"
RELEASE_URL="https://api.github.com/repos/cchalop1/JustDeploy/releases/latest"
SUDO=""

# ——————————————————————————————————————
# Fonction : élévation des privilèges (sudo une fois)
# ——————————————————————————————————————
elevate_priv() {
  if [ "$EUID" -ne 0 ]; then
    if ! command -v sudo &> /dev/null; then
      echo "🛑 sudo n'est pas installé : exécutez le script en root." >&2
      exit 1
    fi
    sudo -v >&2 || {
      echo "🛑 Impossible d'obtenir les privilèges sudo." >&2
      exit 1
    }
    SUDO="sudo"
  fi
}

# ——————————————————————————————————————
# Fonction : teste si on peut écrire dans un dossier
# ——————————————————————————————————————
test_writeable() {
  local dir="$1"
  if touch "${dir}/.perm_test" &> /dev/null; then
    rm -f "${dir}/.perm_test"
    return 0
  else
    return 1
  fi
}

# ——————————————————————————————————————
# 1) Vérification et désinstallation éventuelle
# ——————————————————————————————————————
check_existing_installation() {
  echo "🔍 Vérification d'une installation existante…"
  if [ -f "${BIN_DIR}/justdeploy" ]; then
    echo "🔄 JustDeploy déjà installé, préparation de la réinstallation…"
    if systemctl is-active --quiet "${SERVICE_NAME}"; then
      echo "🛑 Arrêt du service ${SERVICE_NAME}…"
      ${SUDO} systemctl stop "${SERVICE_NAME}"
    fi
    if systemctl is-enabled --quiet "${SERVICE_NAME}"; then
      echo "🔧 Désactivation du service…"
      ${SUDO} systemctl disable "${SERVICE_NAME}"
    fi
    echo "🗑️ Suppression de l’ancien binaire…"
    ${SUDO} rm -f "${BIN_DIR}/justdeploy"
    return 0
  else
    echo "🆕 Aucune installation détectée, installation fraîche."
    return 1
  fi
}

# ——————————————————————————————————————
# 2) Installation des prérequis système (unzip…)
# ——————————————————————————————————————
install_prerequisites() {
  echo "🔍 Vérification des prérequis…"
  if ! command -v unzip &> /dev/null; then
    echo "📦 Installation de unzip…"
    ${SUDO} apt-get update
    ${SUDO} apt-get install -y unzip
    echo "✅ unzip installé."
  else
    echo "✅ unzip déjà présent."
  fi
}

# ——————————————————————————————————————
# 3) Installation de nixpacks
# ——————————————————————————————————————
install_nixpacks() {
  echo "🔍 Vérification de nixpacks…"
  if command -v nixpacks &> /dev/null; then
    echo "✅ nixpacks déjà installé."
  else
    echo "📦 Installation de nixpacks…"
    curl -sSL https://nixpacks.com/install.sh | bash
    echo "✅ nixpacks installé."
  fi
}

# ——————————————————————————————————————
# 4) Installation de Docker (Debian/Ubuntu)
# ——————————————————————————————————————
install_docker() {
  echo "🔍 Vérification de Docker…"
  if ! command -v docker &> /dev/null; then
    echo "🐳 Installation de Docker…"
    ${SUDO} apt-get update
    ${SUDO} apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release
    curl -fsSL \
      https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]')/gpg \
      | ${SUDO} gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] \
      https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]') \
      $(lsb_release -cs) stable" \
      | ${SUDO} tee /etc/apt/sources.list.d/docker.list > /dev/null
    ${SUDO} apt-get update
    ${SUDO} apt-get install -y docker-ce docker-ce-cli containerd.io
    ${SUDO} groupadd -f docker
    ${SUDO} usermod -aG docker "$USER"
    ${SUDO} systemctl enable docker
    ${SUDO} systemctl start docker
    echo "✅ Docker installé."
  else
    echo "✅ Docker déjà présent."
  fi
}

# ——————————————————————————————————————
# 5) Installation de Docker Compose
# ——————————————————————————————————————
install_docker_compose() {
  echo "🔍 Vérification de Docker Compose…"
  if ! command -v docker-compose &> /dev/null; then
    echo "🐳 Installation de Docker Compose…"
    local ver
    ver=$(curl -s https://api.github.com/repos/docker/compose/releases/latest \
      | grep '"tag_name"' | cut -d\" -f4)
    ${SUDO} curl -L \
      "https://github.com/docker/compose/releases/download/${ver}/docker-compose-$(uname -s)-$(uname -m)" \
      -o /usr/local/bin/docker-compose
    ${SUDO} chmod +x /usr/local/bin/docker-compose
    echo "✅ Docker Compose installé."
  else
    echo "✅ Docker Compose déjà présent."
  fi
}

# ——————————————————————————————————————
# Début de l’exécution
# ——————————————————————————————————————
elevate_priv
check_existing_installation
is_reinstall=$?

install_prerequisites
install_nixpacks

platform=$(uname -s | tr '[:upper:]' '[:lower:]')
if [ "$platform" != "darwin" ]; then
  install_docker
  install_docker_compose
fi

# ——————————————————————————————————————
# Détermination du zip et du binaire selon OS/arch
# ——————————————————————————————————————
if [ "$platform" = "darwin" ]; then
  if [ "$(uname -m)" = "arm64" ]; then
    zip_name="justdeploy-darwin-arm.zip"
    bin_src="justdeploy-darwin-arm"
  else
    zip_name="justdeploy-darwin-x86.zip"
    bin_src="justdeploy-darwin-x86"
  fi
else
  arch=$(uname -m)
  if [[ "$arch" == arm* ]] || [[ "$arch" == aarch64 ]]; then
    zip_name="justdeploy-linux-arm.zip"
    bin_src="justdeploy-linux-arm"
  else
    zip_name="justdeploy-linux-x86.zip"
    bin_src="justdeploy-linux-x86"
  fi
fi

# ——————————————————————————————————————
# Récupération du lien de téléchargement
# ——————————————————————————————————————
response=$(curl -sSL "$RELEASE_URL")
download_url=$(echo "$response" \
  | grep -o "https://github.com/cchalop1/JustDeploy/releases/download/[^ ]*/${zip_name}" \
  | head -n1)

# ——————————————————————————————————————
# Téléchargement & installation du binaire
# ——————————————————————————————————————
echo "📥 Téléchargement de JustDeploy…"
curl -sSL -o "$zip_name" "$download_url"

echo "📦 Extraction…"
unzip -o "$zip_name" -d ./bin

echo "🔧 Permissions et déplacement…"
chmod +x ./bin/"$bin_src"
${SUDO} mv "./bin/${bin_src}" "${BIN_DIR}/justdeploy"

# nettoyage
rm -f "$zip_name"
rm -rf ./bin

if [ "$is_reinstall" -eq 0 ]; then
  echo "✨ Réinstallation du binaire terminée."
else
  echo "✨ Installation du binaire terminée."
fi

# ——————————————————————————————————————
# Création du service systemd
# ——————————————————————————————————————
echo "🔧 Création du service systemd…"
cat > /tmp/${SERVICE_NAME} << EOF
[Unit]
Description=JustDeploy Service
After=network.target

[Service]
Type=simple
ExecStart=${BIN_DIR}/justdeploy
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

${SUDO} mv /tmp/${SERVICE_NAME} "${SYSTEMD_DIR}/"
${SUDO} systemctl daemon-reload
${SUDO} systemctl enable "${SERVICE_NAME}"
${SUDO} systemctl start "${SERVICE_NAME}"

echo "✅ Service ${SERVICE_NAME} installé et démarré."

# ——————————————————————————————————————
# Affichage des logs de démarrage
# ——————————————————————————————————————
echo ""
echo "📋 Logs de démarrage (dernières 20 lignes) :"
echo "----------------------------------------------------------------"
sleep 3
${SUDO} journalctl -u "${SERVICE_NAME}" -n 20 --no-pager
echo "----------------------------------------------------------------"
echo "💡 Pour suivre en temps réel : sudo journalctl -u ${SERVICE_NAME} -f"

# ——————————————————————————————————————
# Résumé final
# ——————————————————————————————————————
echo ""
echo "🎉 JustDeploy est opérationnel !"
echo "  • Binaire : ${BIN_DIR}/justdeploy"
echo "  • Service : ${SERVICE_NAME}"
if [ "$platform" != "darwin" ]; then
  command -v docker &> /dev/null && echo "  • Docker : installé"
  command -v docker-compose &> /dev/null && echo "  • Docker Compose : installé"
fi
echo "  • unzip : présent"
echo "  • nixpacks : $(command -v nixpacks &> /dev/null && echo 'présent' || echo 'installé')"
echo ""
echo "🚀 Vous pouvez maintenant utiliser JustDeploy !"
