#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Installing Hyperledger Fabric Prerequisites${NC}"
echo "================================================"

# Download Fabric bootstrap script
curl -sSL https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/bootstrap.sh -o bootstrap.sh
chmod +x bootstrap.sh

# Download Fabric binaries and docker images
./bootstrap.sh

# Copy test-network to our project
cp -r fabric-samples/test-network ./fabric-samples/
cp -r fabric-samples/config ./fabric-samples/

# Add Fabric binaries to PATH
export PATH=$PATH:${PWD}/fabric-samples/bin

# Create necessary directories if they don't exist
mkdir -p network/organizations

echo -e "\n${GREEN}Installation completed!${NC}"
echo "================================================"
echo -e "Please run: ${GREEN}source ~/.bashrc${NC} or ${GREEN}source ~/.zshrc${NC} to update your PATH"
echo -e "Then cd into fabric-samples/test-network to start the network"
