#!/bin/bash

# This script automates the setup of the Hyperledger Fabric backend.
# It should be run from the root of the project directory.

# --- 1. Check for fabric-samples repository ---
echo "--- Checking for fabric-samples repository ---"
if [ -d "../fabric-samples" ]; then
    echo "fabric-samples repository found."
else
    echo "fabric-samples repository not found."
    echo "Please clone it by running the following command in the parent directory of this project:"
    echo "git clone https://github.com/hyperledger/fabric-samples.git"
    exit 1
fi

# --- 2. Bring up the Hyperledger Fabric test network ---
echo "--- Bringing up the Hyperledger Fabric test network ---"
cd ../fabric-samples/test-network
./network.sh down
./network.sh up createChannel -ca -s couchdb

# --- 3. Deploy the smart contract ---
echo "--- Deploying the smart contract ---"
./network.sh deployCC -ccn asset-tracker -ccp ../../Hyperledger-Fabric-Asset-Tracker-Go/chaincode-go -ccl go

# --- 4. Instructions for rest-api-go configuration ---
echo "--- Backend setup complete! ---"
echo ""
echo "--- Instructions for rest-api-go ---"
echo "The REST API needs environment variables to connect to the Fabric network."
echo "The API is currently hardcoded to use the User1@org1.example.com identity."
echo "For a production setup, you should externalize these configurations."
echo ""
echo "The API reads crypto material from a path that should be mounted into its container."
echo "The required crypto material is located in your 'fabric-samples/test-network/organizations' directory."
echo ""

# --- 5. Instructions to run the REST API ---
echo "--- Running the rest-api-go application ---"
echo "To build and run the REST API, navigate to the 'rest-api-go' directory"
echo "in a new terminal and execute the following commands:"
echo ""
echo "1. Build the Docker image:"
echo "   cd ../../Hyperledger-Fabric-Asset-Tracker-Go/rest-api-go"
echo "   docker build -t asset-tracker-api ."
echo ""
echo "2. Run the Docker container, mounting the crypto material:"
echo "   docker run -p 8080:8080 \\"
echo "     -v \`pwd\`/../../fabric-samples/test-network/organizations:/path/to/crypto/material \\"
echo "     asset-tracker-api"
echo ""
echo "Note: Make sure you are in the 'rest-api-go' directory when running the docker run command,"
echo "or adjust the volume path accordingly."
echo ""
echo "Once the container is running, the API will be accessible at http://localhost:8080"
