#!/bin/bash

# Define the release URL and the binary file name
release_url="https://api.github.com/repos/cchalop1/JustDeploy/releases/latest"
zip_file="justdeploy.zip"
binary_file="justdeploy"

# Get the current platform and architecture
platform=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)

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

# Create a log directory
mkdir -p /var/log/justdeploy

# Run the binary in the background and redirect output to log file
echo "🚀 Starting JustDeploy in the background..."
nohup /usr/local/bin/$binary_file > /var/log/justdeploy/justdeploy.log 2>&1 &

# Store the process ID
PID=$!
echo "✅ JustDeploy started with PID: $PID"
echo "📝 Logs are available at: /var/log/justdeploy/justdeploy.log"
echo "💡 You can view logs in real-time with: tail -f /var/log/justdeploy/justdeploy.log"

tail -f /var/log/justdeploy/justdeploy.log

