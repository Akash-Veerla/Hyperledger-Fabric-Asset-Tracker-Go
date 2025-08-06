# Asset Tracker REST API

## Overview

This REST API provides a way to interact with the `asset-tracker` Hyperledger Fabric smart contract. It allows you to create, read, update, delete, and list assets on the blockchain.

## Prerequisites

- Docker
- A running Hyperledger Fabric network with the `asset-tracker` chaincode deployed.
- The crypto material for a user who can interact with the chaincode.

## Configuration

The API is configured using the following environment variables. These are hardcoded in the `main.go` file for now, but should be externalized for a production environment.

- `CRYPTO_PATH`: The path to the directory containing the crypto material. This directory should be mounted into the container.
- `CERT_PATH`: The path to the user's certificate file.
- `KEY_PATH`: The path to the user's private key file.
- `TLS_CERT_PATH`: The path to the TLS CA certificate file for the peer.
- `PEER_ENDPOINT`: The endpoint of the Fabric peer (e.g., `localhost:7051`).
- `GATEWAY_PEER`: The name of the gateway peer (e.g., `peer0.org1.example.com`).
- `CHANNEL_NAME`: The name of the channel the chaincode is deployed on.
- `CHAINCODE_NAME`: The name of the chaincode.

## Running the API

1. **Build the Docker image:**
   ```bash
   docker build -t asset-tracker-api .
   ```

2. **Run the Docker container:**
   ```bash
   docker run -p 8080:8080 -v /path/to/your/crypto/material:/path/to/crypto/material asset-tracker-api
   ```
   Replace `/path/to/your/crypto/material` with the actual path to your crypto material on the host machine.

## API Endpoints

- `POST /assets`: Create a new asset.
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
- `GET /assets/{id}`: Read a specific asset.
- `PUT /assets/{id}`: Update an existing asset.
- `DELETE /assets/{id}`: Delete an asset.
- `GET /assets`: Retrieve all assets.
