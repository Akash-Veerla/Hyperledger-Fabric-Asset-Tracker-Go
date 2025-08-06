# Asset Tracker REST API

## Overview

This REST API provides a way to interact with the `asset-tracker` Hyperledger Fabric smart contract. It allows you to create, read, update, delete, and list assets on the blockchain.

## Prerequisites

- Docker
- A running Hyperledger Fabric network with the `asset-tracker` chaincode deployed.
- The crypto material for a user who can interact with the chaincode.

## Setting up the Hyperledger Fabric Network

This guide assumes you have a basic understanding of Hyperledger Fabric. For a detailed guide on setting up a Fabric network, please refer to the [official Hyperledger Fabric documentation](https://hyperledger-fabric.readthedocs.io/en/latest/test_network.html).

A sample test network is provided in the official Fabric samples. You can use the following steps to set up a test network:

1. **Clone the fabric-samples repository:**
   ```bash
   git clone https://github.com/hyperledger/fabric-samples.git
   ```
2. **Navigate to the test-network directory:**
   ```bash
   cd fabric-samples/test-network
   ```
3. **Start the network:**
   ```bash
   ./network.sh up createChannel
   ```
4. **Deploy the chaincode:**
   ```bash
   ./network.sh deployCC -ccn asset-tracker -ccp ../../chaincode-go -ccl go
   ```

## Generating Crypto Material

The REST API requires crypto material (certificates and private keys) to connect to the Fabric network. The `fabric-samples/test-network` script generates this material for you. You can find it in the `organizations` directory within the `test-network` directory.

The API is configured to use the `User1@org1.example.com` identity. You will need to mount the `organizations` directory into the Docker container when running the API.

## Configuration

The API is configured using the following environment variables. These are hardcoded in the `main.go` file for now, but should be externalized for a production environment.

- `CRYPTO_PATH`: The path to the directory containing the crypto material. This directory should be mounted into the container. For the test network, this would be `.../fabric-samples/test-network/organizations`.
- `CERT_PATH`: The path to the user's certificate file. Default: `CRYPTO_PATH/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/cert.pem`
- `KEY_PATH`: The path to the user's private key file. Default: `CRYPTO_PATH/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/`
- `TLS_CERT_PATH`: The path to the TLS CA certificate file for the peer. Default: `CRYPTO_PATH/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt`
- `PEER_ENDPOINT`: The endpoint of the Fabric peer (e.g., `localhost:7051`).
- `GATEWAY_PEER`: The name of the gateway peer (e.g., `peer0.org1.example.com`).
- `CHANNEL_NAME`: The name of the channel the chaincode is deployed on (e.g., `mychannel`).
- `CHAINCODE_NAME`: The name of the chaincode (e.g., `asset-tracker`).

## Running the API

1. **Build the Docker image:**
   ```bash
   docker build -t asset-tracker-api .
   ```

2. **Run the Docker container:**
   ```bash
   docker run -p 8080:8080 -v /path/to/your/fabric-samples/test-network/organizations:/path/to/crypto/material asset-tracker-api
   ```
   Replace `/path/to/your/fabric-samples/test-network/organizations` with the actual path to your `organizations` directory on the host machine.

## API Endpoints

### POST /assets

Create a new asset.

- **Request body:**
  ```json
  {
    "DEALERID": "dealer3",
    "MSISDN": "111222333",
    "MPIN": "5678",
    "BALANCE": "3000",
    "STATUS": "active",
    "TRANSAMOUNT": "0",
    "TRANSTYPE": "init",
    "REMARKS": "New asset"
  }
  ```
- **Response:**
  - `201 Created`: The asset was created successfully.
  - `400 Bad Request`: The request body is invalid.
  - `500 Internal Server Error`: An error occurred on the server.

### GET /assets/{id}

Read a specific asset.

- **URL parameters:**
  - `id`: The ID of the asset to read.
- **Response:**
  - `200 OK`: The asset was found.
  - `404 Not Found`: The asset was not found.
  - `500 Internal Server Error`: An error occurred on the server.

### PUT /assets/{id}

Update an existing asset.

- **URL parameters:**
  - `id`: The ID of the asset to update.
- **Request body:**
  ```json
  {
    "DEALERID": "dealer3",
    "MSISDN": "111222333",
    "MPIN": "5678",
    "BALANCE": "4000",
    "STATUS": "active",
    "TRANSAMOUNT": "1000",
    "TRANSTYPE": "credit",
    "REMARKS": "Updated asset"
  }
  ```
- **Response:**
  - `200 OK`: The asset was updated successfully.
  - `400 Bad Request`: The request body is invalid.
  - `404 Not Found`: The asset was not found.
  - `500 Internal Server Error`: An error occurred on the server.

### DELETE /assets/{id}

Delete an asset.

- **URL parameters:**
  - `id`: The ID of the asset to delete.
- **Response:**
  - `204 No Content`: The asset was deleted successfully.
  - `404 Not Found`: The asset was not found.
  - `500 Internal Server Error`: An error occurred on the server.

### GET /assets

Retrieve all assets.

- **Response:**
  - `200 OK`: A list of all assets.
  - `500 Internal Server Error`: An error occurred on the server.

## Customization

### Chaincode

The chaincode can be customized by editing the `chaincode-go/asset_tracker.go` file. You can add new fields to the `Asset` struct, add new functions to the smart contract, or modify the existing functions.

After modifying the chaincode, you will need to redeploy it to the Fabric network.

### REST API

The REST API can be customized by editing the `rest-api-go/main.go` file. You can add new endpoints to the API, modify the existing endpoints, or change the way the API interacts with the chaincode.

After modifying the API, you will need to rebuild the Docker image and restart the container.
