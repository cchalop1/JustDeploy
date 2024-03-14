#!/bin/bash

# Define the release URL and the binary file name
release_url="https://api.github.com/repos/cchalop1/JustDeploy/releases/latest"
zip_file="justdeploy.zip"
binary_file="justdeploy"

# Get the latest release download URL for the specific platform
response=$(curl -s $release_url)
download_url=$(echo $response | grep -o "https://github.com/cchalop1/JustDeploy/releases/download/[^ ]*/$zip_file" | head -n 1)

# Download the binary
curl -L -o $zip_file $download_url

# Unzip binary file
unzip $zip_file

# Make the binary executable
chmod +x bin/$binary_file

# Move the binary to a system directory (e.g., /usr/local/bin)
sudo mv bin/$binary_file /usr/local/bin/

rm $zip_file

echo "✨ Installation complete. You can now run $binary_file"