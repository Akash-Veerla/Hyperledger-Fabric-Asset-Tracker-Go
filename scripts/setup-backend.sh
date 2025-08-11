#!/bin/bash

# This script automates the setup of the Hyperledger Fabric backend.
# It should be run from the root of the project directory.

# --- 1. Check for required tools ---
echo "--- Checking for required tools (Docker, Go) ---"
command -v docker >/dev/null 2>&1 || { echo >&2 "Docker is not installed. Aborting."; exit 1; }
command -v go >/dev/null 2>&1 || { echo >&2 "Go is not installed. Aborting."; exit 1; }
echo "All required tools are installed."

# --- 2. Check for fabric-samples repository ---
echo "--- Checking for fabric-samples repository ---"
FABRIC_SAMPLES_PATH="../fabric-samples"
if [ -d "$FABRIC_SAMPLES_PATH" ]; then
    echo "fabric-samples repository found at $FABRIC_SAMPLES_PATH"
else
    echo "fabric-samples repository not found."
    echo "Please clone it by running the following command in the parent directory of this project:"
    echo "git clone https://github.com/hyperledger/fabric-samples.git"
    exit 1
fi

# --- 3. Bring up the Hyperledger Fabric test network ---
echo "--- Bringing up the Hyperledger Fabric test network ---"
cd "$FABRIC_SAMPLES_PATH/test-network"
./network.sh down
./network.sh up createChannel -ca -s couchdb

# --- 4. Deploy the smart contract ---
echo "--- Deploying the smart contract ---"
# Note: The path to the chaincode is relative to the test-network directory
./network.sh deployCC -ccn asset-tracker -ccp ../../Hyperledger-Fabric-Asset-Tracker-Go/chaincode-go -ccl go

# --- 5. Generate config.json for the REST API ---
echo "--- Generating config.json for the REST API ---"
# This path will be the one used inside the rest-api container
CRYPTO_PATH_IN_CONTAINER="/crypto/organizations"
API_CONFIG_PATH="../../Hyperledger-Fabric-Asset-Tracker-Go/rest-api-go/config.json"

# Create the config.json file
cat > "$API_CONFIG_PATH" <<EOL
{
  "mspID": "Org1MSP",
  "cryptoPath": "${CRYPTO_PATH_IN_CONTAINER}",
  "certPath": "peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/cert.pem",
  "keyPath": "peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/",
  "tlsCertPath": "peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt",
  "peerEndpoint": "172.17.0.1:7051",
  "gatewayPeer": "peer0.org1.example.com",
  "channelName": "mychannel",
  "chaincodeName": "asset-tracker",
  "port": "8080"
}
EOL
echo "config.json created at $API_CONFIG_PATH"
# Note: The peerEndpoint is set to the default Docker bridge IP. This might need to be adjusted depending on the user's Docker network setup.

# --- 6. Final Instructions ---
echo ""
echo "--- Backend setup complete! ---"
echo ""
echo "The Hyperledger Fabric network is running, and the smart contract is deployed."
echo "The REST API configuration has been generated."
echo ""
echo "You can now run the entire application stack using Docker Compose."
echo "Navigate to the root of the project directory and run:"
echo ""
echo "docker-compose up --build"
echo ""
