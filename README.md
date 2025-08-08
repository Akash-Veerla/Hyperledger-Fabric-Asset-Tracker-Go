# Hyperledger Fabric Asset Tracker

This project is a blockchain-based system for managing and tracking financial assets, built on Hyperledger Fabric. It includes a Go smart contract, a REST API, and a React front-end.

## Project Overview

The project is divided into three main components:

1.  **`chaincode-go`**: This directory contains the Go smart contract (chaincode) that runs on the Hyperledger Fabric network. The chaincode is responsible for managing the assets on the blockchain, including creating, reading, updating, and deleting assets.

2.  **`rest-api-go`**: This directory contains a REST API built with Go and the Gin framework. The API acts as a bridge between the front-end and the blockchain network, allowing you to interact with the smart contract through standard HTTP requests.

3.  **`frontend-react`**: This directory contains a single-page application (SPA) built with React and Tailwind CSS. The front-end provides a user-friendly interface for managing the financial assets, interacting with the REST API to perform operations on the blockchain.

## Getting Started

To get started with this project, you will need to set up the Hyperledger Fabric network, deploy the smart contract, run the REST API, and run the front-end application.

### 1. Set up the Hyperledger Fabric Network

Follow the instructions in the `rest-api-go/README.md` file to set up a local Hyperledger Fabric network using the `fabric-samples`.

### 2. Deploy the Smart Contract

Once you have the Fabric network running, deploy the smart contract using the following command from the `fabric-samples/test-network` directory:

```bash
./network.sh deployCC -ccn asset-tracker -ccp ../../chaincode-go -ccl go
```

### 3. Run the REST API

Follow the instructions in the `rest-api-go/README.md` file to build and run the REST API in a Docker container.

### 4. Run the Front-End Application

Follow the instructions in the `frontend-react/README.md` file to set up and run the React front-end application.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
