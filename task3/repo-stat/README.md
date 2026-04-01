# GitHub Repository Info Tool

This project is a tool for getting information about a GitHub repository.

## Features

- Fetch repository data by repository URL
- REST API via `api` gateway
- gRPC services: `processor`, `collector`, `subscriber`
- Swagger web interface for API testing

## Usage

1. Go to the project root directory:

   ```bash
   cd ..
   ```

2. Start the services:

   ```bash
   make up
   ```

3. Open Swagger UI in your browser:

   `http://localhost:28080/swagger/index.html`

4. Run integration tests (optional):

   ```bash
   make test
   ```

5. Stop and remove containers when you are done:

   ```bash
   make down
   ```
