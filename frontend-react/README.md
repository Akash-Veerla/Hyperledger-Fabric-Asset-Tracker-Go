# Asset Tracker Front-End

This is a single-page application (SPA) built with React and Vite, providing a user interface for the Asset Tracker project.

## Prerequisites

- Node.js and npm
- A running instance of the `rest-api-go` service.

## Getting Started

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Run the development server:**
   ```bash
   npm run dev
   ```

This will start the Vite development server, and you can view the application in your browser at `http://localhost:5173` (or the address shown in your terminal).

## Connecting to the REST API

The front-end is configured to connect to the REST API at `http://localhost:8080`. If your API is running on a different address, you will need to update the `API_URL` constant in the `src/App.jsx` file.
