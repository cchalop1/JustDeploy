#!/bin/bash

# Define the release URL and the binary file name
release_url="https://api.github.com/repos/cchalop1/JustDeploy/releases/latest"
binary_file="justdeploy.zip"

# Get the latest release download URL for the specific platform
response=$(curl -s $release_url)
download_url=$(echo $response | grep -o "https://github.com/your_username/your_repo/releases/download/[^ ]*/$binary_file" | head -n 1)

# Download the binary
curl -L -o $binary_file $download_url

# Make the binary executable
chmod +x $binary_file

# Move the binary to a system directory (e.g., /usr/local/bin)
sudo mv $binary_file /usr/local/bin/

echo "Installation complete. You can now run $binary_file"