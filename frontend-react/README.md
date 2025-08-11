# Asset Tracker Front-End

## Overview

This is a single-page application (SPA) built with React and Vite, providing a user interface for the Asset Tracker project. It is intended to be run as part of the `docker-compose` setup in the root of the project.

## Running the Front-End

The recommended way to run this application is through the main `docker-compose.yaml` file in the root of the project.

```bash
# From the project root
docker-compose up --build frontend
```

This will build the Docker image and start the container, serving the application on port 3000. The service is also started as part of the full stack `docker-compose up` command.

### Development

For local development, you can run the Vite dev server:

1.  **Install dependencies:**
    ```bash
    npm install
    ```

2.  **Run the development server:**
    ```bash
    npm run dev
    ```

The application will be available at `http://localhost:5173`.
