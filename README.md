# Hyperledger Fabric Asset Tracker

This project is a blockchain-based system for managing and tracking financial assets, built on Hyperledger Fabric. It includes a Go smart contract, a REST API, and a React front-end.

## Architecture

The project consists of three main services orchestrated with Docker Compose:
1.  **`chaincode-go`**: The Go smart contract that runs on the Hyperledger Fabric network. It defines the business logic for managing assets.
2.  **`rest-api-go`**: A REST API built with Go that acts as a bridge to the blockchain network.
3.  **`frontend-react`**: A single-page application built with React and Vite that provides a user interface for interacting with the system.

## Getting Started

The entire application stack can be set up and run with a combination of a setup script and Docker Compose.

### Prerequisites
- Docker and Docker Compose
- Go
- Node.js and npm
- A local clone of the `fabric-samples` repository in the parent directory of this project.

### 1. Set up the Backend

First, run the backend setup script. This script will:
- Start the Hyperledger Fabric test network.
- Deploy the smart contract.
- Generate the necessary configuration for the REST API.

From the root of this project, run:
```bash
cd scripts
./setup-backend.sh
```

### 2. Run the Application Stack

Once the backend setup is complete, you can start the entire application using Docker Compose.

From the root of this project, run:
```bash
docker-compose up --build
```

This will:
- Build the Docker images for the REST API and the frontend.
- Start the containers.

You can then access:
- The frontend application at `http://localhost:3000`
- The REST API at `http://localhost:8080`

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
