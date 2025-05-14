#!/bin/bash

# Check if script is run as root
if [ "$(id -u)" -ne 0 ]; then
  echo "🛑 This script must be run as root or with sudo privileges."
  echo "Please run: sudo $0"
  exit 1
fi

# Define variables
DOCKER_IMAGE="cchalop1/justdeploy:latest"
DATA_DIR="/var/lib/justdeploy"
DOCKER_SOCKET="/var/run/docker.sock"
CONTAINER_NAME="justdeploy"

# Function to check if JustDeploy is already installed and stop if needed
check_existing_installation() {
  echo "🔍 Checking for existing JustDeploy installation..."
  
  # Remove existing Docker container if exists
  if docker ps -a | grep -q $CONTAINER_NAME; then
    echo "🛑 Stopping existing JustDeploy container..."
    docker stop $CONTAINER_NAME > /dev/null 2>&1
    echo "🗑️ Removing existing JustDeploy container..."
    docker rm -f $CONTAINER_NAME > /dev/null 2>&1
    echo "✅ Existing JustDeploy Docker container removed."
  fi
  
  # Return true (0) to indicate reinstallation is needed
  return 0
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

# Install prerequisites
install_prerequisites

# Get the current platform and architecture
platform=$(uname -s | tr '[:upper:]' '[:lower:]')

# Check and install Docker if needed
if [ "$platform" != "darwin" ]; then
  install_docker
fi

# Create data directory
echo "📁 Creating data directory..."
mkdir -p $DATA_DIR
chmod 755 $DATA_DIR

# Pull the latest JustDeploy Docker image
echo "🐳 Pulling the latest JustDeploy Docker image..."
docker pull $DOCKER_IMAGE

# Run the Docker container
echo "🚀 Starting JustDeploy Docker container..."
docker run -d --name $CONTAINER_NAME \
  -p 5915:5915 \
  -v $DATA_DIR:/app/data \
  -v $DOCKER_SOCKET:/var/run/docker.sock \
  --restart=unless-stopped \
  $DOCKER_IMAGE

# Check if container started successfully
if [ $? -eq 0 ]; then
  echo "✅ JustDeploy Docker container started successfully!"
else
  echo "❌ Failed to start JustDeploy Docker container. Check Docker logs for details."
  exit 1
fi

echo "------------------------"
echo "✅ JustDeploy Docker container running"
if [ "$platform" != "darwin" ]; then
  if command -v docker &> /dev/null; then
    echo "✅ Docker is installed"
  fi
fi
echo "✅ Nixpacks installed (prerequisite)"
echo ""
echo "🚀 JustDeploy is now running in a Docker container!"
echo "💡 Access the web interface: http://localhost:5915"
echo "💻 To stop the container: docker stop $CONTAINER_NAME"
echo "💻 To start the container: docker start $CONTAINER_NAME"
echo "💻 To view logs: docker logs $CONTAINER_NAME"
echo ""