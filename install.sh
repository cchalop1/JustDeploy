#!/bin/bash

# Check if script is run as root
if [ "$(id -u)" -ne 0 ]; then
  echo "🛑 This script must be run as root or with sudo privileges."
  echo "Please run: sudo $0"
  exit 1
fi

# Define the release URL and the binary file name
release_url="https://api.github.com/repos/cchalop1/JustDeploy/releases/latest"
zip_file="justdeploy.zip"
binary_file="justdeploy"

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

# Get the current platform and architecture
platform=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

# Check and install Docker and Docker Compose if needed (only on Linux)
if [ "$platform" != "darwin" ]; then
  install_docker
  install_docker_compose
fi

if [ "$platform" == "darwin" ]; then
  if [ "$arch" == "arm64" ]; then
    zip_file="justdeploy-darwin-arm.zip"
    binary_file_arch="justdeploy-darwin-arm"
  else
    zip_file="justdeploy-darwin-x86.zip"
    binary_file_arch="justdeploy-darwin-x86"
  fi
else
  if [ "$(expr substr $(uname -m) 1 5)" == "armv7" ] || [ "$(expr substr $(uname -m) 1 3)" == "aarch64" ]; then
    zip_file="justdeploy-linux-arm.zip"
    binary_file_arch="justdeploy-linux-arm"
  else
    zip_file="justdeploy-linux-x86.zip"
    binary_file_arch="justdeploy-linux-x86"
  fi
fi

# Get the latest release download URL for the specific platform
response=$(curl -s $release_url)
download_url=$(echo $response | grep -o "https://github.com/cchalop1/JustDeploy/releases/download/[^ ]*/$zip_file" | head -n 1)

# Download the binary
curl -L -o $zip_file $download_url

# Unzip binary file
unzip $zip_file

# Make the binary executable
chmod +x ./bin/$binary_file_arch

# Move the binary to a system directory (e.g., /usr/local/bin)
sudo mv ./bin/$binary_file_arch /usr/local/bin/$binary_file

rm $zip_file

rm -rf ./bin

echo "✨ Installation complete. You can now run $binary_file"

# # Create systemd service file
# echo "🔧 Creating systemd service for JustDeploy..."
# cat > /tmp/justdeploy.service << EOF
# [Unit]
# Description=JustDeploy Service
# After=network.target

# [Service]
# Type=simple
# ExecStart=/usr/local/bin/$binary_file
# Restart=on-failure
# RestartSec=10
# StandardOutput=journal
# StandardError=journal

# [Install]
# WantedBy=multi-user.target
# EOF

# # Move service file to systemd directory
# sudo mv /tmp/justdeploy.service /etc/systemd/system/

# # Reload systemd to recognize the new service
# sudo systemctl daemon-reload

# # Enable and start the service
# sudo systemctl enable justdeploy.service
# sudo systemctl start justdeploy.service

# echo "✅ JustDeploy service has been installed and started"
# echo "💡 You can check the status with: sudo systemctl status justdeploy.service"
# echo "📝 View logs with: sudo journalctl -u justdeploy.service -f"

# # Show the current status
# sudo systemctl status justdeploy.service

# # Print summary
# echo ""
# echo "🎉 Installation Summary:"
# echo "------------------------"
# echo "✅ JustDeploy installed at: /usr/local/bin/$binary_file"
# echo "✅ Systemd service created: justdeploy.service"
# if [ "$platform" != "darwin" ]; then
#   if command -v docker &> /dev/null; then
#     echo "✅ Docker is installed"
#   fi
#   if command -v docker-compose &> /dev/null; then
#     echo "✅ Docker Compose is installed"
#   fi
# fi
# echo ""
# echo "🚀 JustDeploy is now running as a system service!"
# echo "💡 Access the web interface using the URL shown in the service logs"
# echo "📝 View logs with: sudo journalctl -u justdeploy.service -f"
# echo ""

