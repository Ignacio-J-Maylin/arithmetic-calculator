
# Arithmetic Calculator API

This project implements an API for arithmetic calculations and user management with authentication. It is developed in **Go** and uses **MySQL** as the database.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Run the Service Locally](#run-the-service-locally)
3. [Build the Docker Image](#build-the-docker-image)
4. [Run the Service with Docker Compose](#run-the-service-with-docker-compose)
5. [Postman Collection](#postman-collection)
6. [Endpoint Description](#endpoint-description)

---

## Prerequisites

Ensure you have the following installed on your local environment:

- [Go](https://golang.org/) (version 1.20 or higher)
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
- [Postman](https://www.postman.com/)

---

## Run the Service Locally

### Configuration

1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd <repo-directory>
   ```

2. Configure environment variables in a `.env` file at the project root:
   ```dotenv
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=admin
   DB_PASSWORD=admin
   DB_NAME=arithmetic
   API_PORT=8080
   ```

3. Start a local MySQL database or use Docker (optional):
   ```bash
   docker run -d      --name arithmetic-db      -e MYSQL_DATABASE=arithmetic      -e MYSQL_USER=admin      -e MYSQL_PASSWORD=admin      -e MYSQL_ROOT_PASSWORD=admin      -p 3306:3306      mysql:8.0
   ```

### Execution

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

The service will be available at `http://localhost:8080`.

---

## Build the Docker Image

1. Build the Docker image with a tag:
   ```bash
   docker build -t arithmetic-calculator:v1 .
   ```

2. Run the container:
   ```bash
   docker run -d -p 8080:8080 --name arithmetic-backend arithmetic-calculator:v1
   ```

---

## Run the Service with Docker Compose

1. Verify the `docker-compose.yaml` file and ensure the configuration matches your environment.
2. Start the services:
   ```bash
   docker-compose up -d
   ```

This will start the backend, frontend, and MySQL database. The backend will be available at `http://localhost:8080`, and the frontend will be available at `http://localhost:3000`.

---

## Postman Collection

### Import the Collection

1. Copy the Postman collection JSON content (provided in the project) and save it as `postman_collection.json`.
2. Open Postman and select **Import** from the top menu.
3. Drag and drop the `postman_collection.json` file or copy and paste its content.
4. Set the following global variables in Postman:
   - `{{token}}`: Token generated after login.
   - `{{refresh_token}}`: Token to refresh the session.

### Using the Collection

1. Test the login endpoint:
   - Method: `POST`
   - URL: `http://localhost:8080/api/v1/login`
   - Body:
     ```json
     {
       "username": "testuser@gmail.com",
       "password": "password123"
     }
     ```

   The authentication token will be automatically set in Postman's global variables.

2. Test other endpoints in the collection, such as `add_credits`, `get_credits`, and arithmetic operations.

---

## Endpoint Description

1. **Authentication**:
   - `POST /api/v1/login`: Logs in and returns an access token.
   - `POST /api/v1/signup`: Registers a new user.
   - `POST /api/v1/logout`: Logs out.

2. **Credit Management**:
   - `PUT /api/v1/users/credits`: Adds or removes credits.
   - `GET /api/v1/users/credits`: Gets the credit balance.

3. **Arithmetic Operations**:
   - `POST /api/v1/users/operation`: Performs operations such as addition, subtraction, multiplication, division, etc.

4. **Record History**:
   - `GET /api/v1/records/history`: Fetches the operation history.
   - `DELETE /api/v1/records/delete`: Deletes a specific record.

---

This README should serve as a comprehensive guide for any developer or tester working with your Go backend project.
