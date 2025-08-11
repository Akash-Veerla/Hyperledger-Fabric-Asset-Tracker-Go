# Hyperledger Fabric Asset Tracker

This project is a blockchain-based system for managing and tracking financial assets, built on Hyperledger Fabric. It includes a Go smart contract, a REST API, and a React front-end.

## Project Overview

The project is divided into three main components:

1.  **`chaincode-go`**: This directory contains the Go smart contract (chaincode) that runs on the Hyperledger Fabric network. The chaincode is responsible for managing the assets on the blockchain, including creating, reading, updating, and deleting assets.

2.  **`rest-api-go`**: This directory contains a REST API built with Go and the Gin framework. The API acts as a bridge between the front-end and the blockchain network, allowing you to interact with the smart contract through standard HTTP requests.

3.  **`frontend-react`**: This directory contains a single-page application (SPA) built with React and Vite. The front-end provides a user-friendly interface for managing the financial assets, interacting with the REST API to perform operations on the blockchain.

4.  **`scripts`**: This directory contains helper scripts to automate the setup and deployment of the project.

## Getting Started

To get started with this project, you can use the automated setup script for the backend, and then run the front-end application.

### 1. Set up the Backend (Fabric Network and REST API)

The `setup-backend.sh` script automates the process of setting up the Hyperledger Fabric network, deploying the smart contract, and provides instructions for running the REST API.

To run the script, navigate to the `scripts` directory and execute the following command:

```bash
./setup-backend.sh
```

The script will guide you through the process.

### 2. Run the Front-End Application

Follow the instructions in the `frontend-react/README.md` file to set up and run the React front-end application.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
