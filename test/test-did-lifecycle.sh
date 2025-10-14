#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# API endpoint
API_URL="http://localhost:3000/api"

echo -e "${GREEN}🔧 Testing DID Registry Lifecycle${NC}"
echo "================================================"

# Test 1: Create a new DID
echo -e "\n${GREEN}1. Creating new DID for IoT device...${NC}"
CREATE_RESPONSE=$(curl -s -X POST $API_URL/did/create \
  -H "Content-Type: application/json" \
  -d '{
    "didID": "did:hflc:sensor-101",
    "pubKey": "ABCDEF123456ABCDEF123456"
  }')

echo "Response: $CREATE_RESPONSE"

# Sleep to allow transaction to be committed
sleep 2

# Test 2: Resolve the DID
echo -e "\n${GREEN}2. Resolving DID to verify creation...${NC}"
RESOLVE_RESPONSE=$(curl -s $API_URL/did/resolve/did:hflc:sensor-101)
echo "Response: $RESOLVE_RESPONSE"

# Test 3: Verify the device's identity
echo -e "\n${GREEN}3. Verifying device identity...${NC}"
if echo $RESOLVE_RESPONSE | grep -q "ABCDEF123456"; then
    echo -e "${GREEN}✓ Device public key verified successfully${NC}"
else
    echo -e "${RED}✗ Device verification failed${NC}"
fi

# Test 4: Revoke the DID
echo -e "\n${GREEN}4. Revoking DID...${NC}"
REVOKE_RESPONSE=$(curl -s -X POST $API_URL/did/revoke \
  -H "Content-Type: application/json" \
  -d '{
    "didID": "did:hflc:sensor-101"
  }')

echo "Response: $REVOKE_RESPONSE"

# Sleep to allow transaction to be committed
sleep 2

# Test 5: Attempt to resolve revoked DID
echo -e "\n${GREEN}5. Attempting to resolve revoked DID...${NC}"
RESOLVE_REVOKED_RESPONSE=$(curl -s $API_URL/did/resolve/did:hflc:sensor-101)
echo "Response: $RESOLVE_REVOKED_RESPONSE"

# Verify revocation status
if echo $RESOLVE_REVOKED_RESPONSE | grep -q "revoked"; then
    echo -e "${GREEN}✓ DID revocation verified successfully${NC}"
else
    echo -e "${RED}✗ DID revocation verification failed${NC}"
fi

echo -e "\n${GREEN}Test suite completed!${NC}"
echo "================================================"
