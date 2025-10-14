# IoT Decentralized Identity (DID) Registry

A Zero-Trust Access Control Framework for Heterogeneous IoT Devices using Decentralized Identifiers (DIDs) on Hyperledger Fabric.

## Project Overview

This project implements a secure, decentralized identity layer for IoT devices using Hyperledger Fabric. It provides a robust framework for managing device identities through W3C Decentralized Identifiers (DIDs) and enables zero-trust verification for device authentication and access control.

## Prerequisites

Before starting, ensure you have the following installed:
* **Docker** and **Docker Compose**
* **Go** (version 1.18 or higher)
* **Node.js** (LTS version) and **npm**
* **curl** for downloading Fabric components

## Setup Instructions

### 1. Install Hyperledger Fabric

First, we need to install Hyperledger Fabric and its prerequisites:

```bash
# Make the installation script executable
chmod +x scripts/install-fabric.sh

# Run the installation script
./scripts/install-fabric.sh

# Update your PATH (use appropriate command for your shell)
source ~/.bashrc  # for bash
# OR
source ~/.zshrc   # for zsh
```

### 2. Deploy the DID Registry Chaincode

```bash
cd fabric-samples/test-network

# Start the network and create channel
./network.sh up createChannel -c didchannel

# Deploy the chaincode
./network.sh deployCC -ccn didregistry -ccp ../../chaincode/did-registry -ccl go -c didchannel
```

### 3. Start the API Gateway

```bash
cd api-gateway
npm install
node src/app.js
```

### 4. Test the DID Registry

```bash
# Make the test script executable
chmod +x test/test-did-lifecycle.sh

# Run the tests
./test/test-did-lifecycle.sh
```

## API Endpoints

### Create DID
```http
POST /api/did/create
Content-Type: application/json

{
    "didID": "did:hflc:sensor-101",
    "pubKey": "ABCDEF123456"
}
```

### Resolve DID
```http
GET /api/did/resolve/did:hflc:sensor-101
```

### Revoke DID
```http
POST /api/did/revoke
Content-Type: application/json

{
    "didID": "did:hflc:sensor-101"
}
```

## Project Structure

```
.
├── api-gateway/              # Node.js API Gateway
│   ├── src/
│   │   ├── app.js           # Express application
│   │   └── fabric-connection.js  # Fabric Gateway connection
│   └── package.json
├── chaincode/               # Hyperledger Fabric Chaincode
│   └── did-registry/
│       ├── did_registry.go  # DID Registry implementation
│       └── go.mod          # Go module file
├── network/                # Network configuration
│   └── organizations/     # Crypto materials
├── scripts/               # Setup scripts
│   └── install-fabric.sh  # Fabric installation script
├── test/                  # Test scripts
│   └── test-did-lifecycle.sh  # DID lifecycle test
└── README.md
```

## Security Considerations

- All DID operations are recorded on the blockchain
- Zero-trust verification model
- Public key cryptography for device authentication
- Permissioned blockchain network

## License

MIT License