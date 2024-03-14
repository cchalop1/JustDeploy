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
    binary_file="justdeploy-darwin-arm"
  else
    zip_file="justdeploy-darwin-x86.zip"
    binary_file="justdeploy-darwin-x86"
  fi
else
  if [ "$(expr substr $(uname -m) 1 5)" == "armv7" ] || [ "$(expr substr $(uname -m) 1 3)" == "aarch64" ]; then
    zip_file="justdeploy-linux-arm.zip"
    binary_file="justdeploy-linux-arm"
  else
    zip_file="justdeploy-linux-x86.zip"
    binary_file="justdeploy-linux-x86"
  fi
fi

# Get the latest release download URL for the specific platform
response=$(curl -s $release_url)
download_url=$(echo $response | grep -o "https://github.com/cchalop1/JustDeploy/releases/download/[^ ]*/$zip_file" | head -n 1)

# Download the binary
curl -L -o $zip_file $download_url

# Unzip binary file
unzip $binary_file $zip_file

# Make the binary executable
chmod +x $binary_file

# Move the binary to a system directory (e.g., /usr/local/bin)
sudo mv $binary_file /usr/local/bin/

rm $zip_file

echo "✨ Installation complete. You can now run $binary_file"