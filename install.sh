#!/bin/bash

# Check if the OS is Linux
if [ "$(uname -s)" != "Linux" ]; then
  echo "🛑 This script is intended for Linux systems only (Debian/Ubuntu)."
  echo "Installation on $(uname -s) is not supported."
  exit 1
fi

# Check if the architecture is supported (x86 or ARM)
arch=$(uname -m)
if [[ "$arch" != "x86_64" && "$arch" != "i686" && "$arch" != "armv7l" && "$arch" != "aarch64" ]]; then
  echo "🛑 Unsupported architecture: $arch"
  echo "This script only supports x86 (32/64-bit) or ARM architectures."
  exit 1
fi

# Check if script is run as root
if [ "$(id -u)" -ne 0 ]; then
  echo "🛑 This script must be run as root or with sudo privileges."
  echo "Please run: sudo $0"
  exit 1
fi

# Work only in a private installation directory; never write predictable root-owned files in /tmp.
install_work_dir=$(mktemp -d) || exit 1
trap 'rm -rf "$install_work_dir"' EXIT
cd "$install_work_dir" || exit 1

# Define the release URL and the binary file name
release_url="https://api.github.com/repos/cchalop1/JustDeploy/releases/latest"
zip_file="justdeploy.zip"
binary_file="justdeploy"

# Function to check if JustDeploy is already installed and stop the service if needed
check_existing_installation() {
  echo "🔍 Checking for existing JustDeploy installation..."
  
  if [ -f "/usr/local/bin/$binary_file" ]; then
    echo "🔄 JustDeploy is already installed. Preparing for reinstallation..."
    
    # Leave the running installation intact until the replacement is downloaded and verified.
    # Return true (0) to indicate reinstallation is needed
    return 0
  else
    echo "🆕 No existing JustDeploy installation detected. Proceeding with fresh installation."
    # Return false (1) to indicate fresh installation
    return 1
  fi
}

# Function to install prerequisites
install_prerequisites() {
  echo "🔍 Checking for required packages..."
  if ! command -v unzip &> /dev/null; then
    echo "📦 Installing unzip..."
    sudo apt-get update
    sudo apt-get install -y unzip
    echo "✅ unzip installed successfully."
  else
    echo "✅ unzip is already installed."
  fi
  
  # Check if Nixpacks is installed
  if ! command -v nixpacks &> /dev/null; then
    echo "📦 Installing Nixpacks..."
    echo "Nixpacks is optional. Install a verified release manually to build repositories without a Dockerfile."
    echo "Dockerfile builds work without Nixpacks."
  else
    echo "✅ Nixpacks is already installed."
  fi
}

# Function to check if Docker is installed and install it if not
install_docker() {
  echo "🔍 Checking if Docker is installed..."
  if command -v docker &> /dev/null; then
    echo "✅ Docker is already installed."
  else
    echo "🐳 Docker not found. Installing Docker..."
    
    # Check if we're on a Debian-based system
    if ! command -v apt-get &> /dev/null; then
      echo "❌ This script only supports Docker installation on Debian-based systems (Ubuntu, Debian, etc.)."
      echo "Please install Docker manually according to your OS instructions: https://docs.docker.com/engine/install/"
      return 1
    fi
    
    # Update package index
    sudo apt-get update
    
    # Install packages to allow apt to use a repository over HTTPS
    sudo apt-get install -y \
      apt-transport-https \
      ca-certificates \
      curl \
      gnupg \
      lsb-release
    
    # Add Docker's official GPG key
    curl -fsSL https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]')/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
    
    # Set up the stable repository
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]') \
      $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    
    # Update apt package index again
    sudo apt-get update
    
    # Install Docker Engine
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io
    
    # Add current user to docker group to run docker without sudo
    sudo groupadd -f docker
    sudo usermod -aG docker $USER
    
    # Enable and start Docker service
    sudo systemctl enable docker
    sudo systemctl start docker
    
    echo "✅ Docker has been installed successfully."
    echo "💡 NOTE: You may need to log out and log back in for the docker group changes to take effect."
  fi
}

# Function to check if Docker Compose is installed and install it if not
install_docker_compose() {
  echo "🔍 Checking if Docker Compose is installed..."
  if command -v docker-compose &> /dev/null; then
    echo "✅ Docker Compose is already installed."
  else
    echo "🐳 Docker Compose not found. Installing Docker Compose..."
    
    # Install Docker Compose
    COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep 'tag_name' | cut -d\" -f4)
    sudo curl -L "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    
    echo "✅ Docker Compose has been installed successfully."
  fi
}

# Check for existing installation first
if check_existing_installation; then is_reinstall=0; else is_reinstall=1; fi

# Install prerequisites
install_prerequisites

# Get the current platform and architecture
platform=$(uname -s | tr '[:upper:]' '[:lower:]')

# Check and install Docker and Docker Compose if needed
install_docker
install_docker_compose

# Determine the appropriate binary based on architecture
if [ "$arch" == "armv7l" ] || [ "$arch" == "aarch64" ]; then
  zip_file="justdeploy-linux-arm.zip"
  binary_file_arch="justdeploy-linux-arm"
else
  zip_file="justdeploy-linux-x86.zip"
  binary_file_arch="justdeploy-linux-x86"
fi

# Get the latest release download URL for the specific platform
response=$(curl -fsSL "$release_url") || exit 1
download_url=$(echo $response | grep -o "https://github.com/cchalop1/JustDeploy/releases/download/[^ ]*/$zip_file" | head -n 1)

# Download the binary
echo "📥 Downloading JustDeploy binary..."
curl -fL --proto '=https' --tlsv1.2 -o "$zip_file" "$download_url" || exit 1
checksum_url="${download_url%/*}/SHA256SUMS"
curl -fL --proto '=https' --tlsv1.2 -o SHA256SUMS "$checksum_url" || exit 1
grep "  $zip_file$" SHA256SUMS | sha256sum -c - || exit 1

# Unzip binary file
echo "📦 Extracting binary..."
unzip -o $zip_file

# Make the binary executable
chmod +x ./bin/$binary_file_arch

# Stop only after verification, immediately before replacing the binary.
if systemctl is-active --quiet justdeploy.service; then systemctl stop justdeploy.service || exit 1; fi

# Move the binary to a system directory (e.g., /usr/local/bin)
sudo install -m 755 "./bin/$binary_file_arch" "/usr/local/bin/$binary_file" || exit 1

# Clean up downloaded files
rm $zip_file
rm -rf ./bin

if [ $is_reinstall -eq 0 ]; then
  echo "✨ JustDeploy binary reinstallation complete."
else
  echo "✨ JustDeploy binary installation complete."
fi

# Create systemd service file
echo "🔧 Creating systemd service for JustDeploy..."
cat > "$install_work_dir/justdeploy.service" << EOF
[Unit]
Description=JustDeploy Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/$binary_file
Environment=JUSTDEPLOY_DATA_DIR=/var/lib/justdeploy
Environment=JUSTDEPLOY_LISTEN_ADDR=127.0.0.1
EnvironmentFile=-/etc/justdeploy.env
UMask=0077
NoNewPrivileges=true
PrivateTmp=true
Restart=on-failure
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# Preserve existing state before switching from the legacy root HOME directory.
if [ -d /root/.config/justdeploy ] && [ ! -e /var/lib/justdeploy ]; then
  sudo cp -a /root/.config/justdeploy /var/lib/justdeploy || exit 1
fi
sudo install -d -m 700 /var/lib/justdeploy

# Move service file to systemd directory
sudo mv "$install_work_dir/justdeploy.service" /etc/systemd/system/

# Reload systemd to recognize the new service
sudo systemctl daemon-reload

# Enable and start the service
sudo systemctl enable justdeploy.service
sudo systemctl start justdeploy.service

if [ $is_reinstall -eq 0 ]; then
  echo "✅ JustDeploy service has been reinstalled and started"
else
  echo "✅ JustDeploy service has been installed and started"
fi

# Display startup logs to show server IP
echo "📋 Displaying JustDeploy startup logs (showing server IP):"
echo "----------------------------------------------------------------"
sleep 3  # Give the service a moment to start
sudo journalctl -u justdeploy.service -n 20 --no-pager
echo "----------------------------------------------------------------"
echo "💡 You can continue to monitor logs with: sudo journalctl -u justdeploy.service -f"

# Print summary
echo ""
if [ $is_reinstall -eq 0 ]; then
  echo "🎉 Reinstallation Summary:"
else
  echo "🎉 Installation Summary:"
fi
echo "------------------------"
echo "✅ JustDeploy installed at: /usr/local/bin/$binary_file"
echo "✅ Systemd service created: justdeploy.service"
if command -v docker &> /dev/null; then
  echo "✅ Docker is installed"
fi
if command -v docker-compose &> /dev/null; then
  echo "✅ Docker Compose is installed"
fi
if command -v unzip &> /dev/null; then
  echo "✅ Unzip installed (prerequisite)"
fi
if command -v nixpacks &> /dev/null; then
  echo "✅ Nixpacks installed (prerequisite)"
fi
echo ""
echo "🚀 JustDeploy is now running as a system service!"
echo "💡 Access the web interface using the URL shown in the service logs above"
echo ""