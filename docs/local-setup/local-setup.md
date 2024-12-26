# Local Environment Documentation (Docker Compose)

## Introduction

This document describes how to set up and run the application in a local environment using Docker Compose.

## Prerequisites

* **Docker:** Installed and running Docker Engine. [Docker installation instructions](https://docs.docker.com/get-docker/)
* **Docker Compose:** Installed and running Docker Compose. [Docker Compose installation instructions](https://docs.docker.com/compose/install/)
* **Make:** Installed Make tool. (Refer to your operating system documentation for installation instructions)

## Project Structure

The project consists of the following services:

* **vault:** Vault server for secrets management.
* **authentication-service:** Authentication service.
* **logger-service:** Logging service.
* **mailer-service:** Email sending service.
* **user-service:** User management service.
* **admin-service:** Administrative service.
* **redis:** Redis server for caching.
* **flyway:** Database migration tool.
* **mailhog:**  Server for testing emails.
* **postgres-users:** PostgreSQL database for users.
* **postgres-admins:** PostgreSQL database for administrators.
* **mongo:** MongoDB database for logs.

## Configuration

The environment is configured using the `docker-compose.yaml` file and `Makefile`.

* **docker-compose.yaml:** Defines services, their dependencies, network configuration, and volumes.
* **Makefile:** Contains commands for building, starting, and stopping containers.

## Running the Environment

1.  **Building images:**

    ```bash
    make up_build
    ```

    This command will build Docker images for all services and start them in the background.

2.  **Starting containers:**

    ```bash
    make up
    ```

    This command will start the containers in the background, without building images.

## Stopping the Environment

1.  **Stopping containers:**

    ```bash
    make down
    ```

    This command will stop the containers.

2.  **Stopping and removing containers:**

    ```bash
    make down_rm
    ```

    This command will stop and remove the containers and their associated data.

## Building the Application

1.  **Building all applications:**

    ```bash
    make build
    ```

    This command will build binaries for all services.

2.  **Building a single application:**

    ```bash
    make build_auth 
    make build_user
    make build_admin
    make build_mail
    make build_logger
    ```

    These commands will build binaries for individual services.

## Testing

After starting the environment, you can test the application by sending requests to individual services.

**Example:**

To test the user-service, you can send a GET request to `http://localhost:8081/users`.

## Additional Information

* **Access to Vault:**
    * Address: `http://localhost:8200`
    * Token: `root_token`
* **Access to Mailhog:** `http://localhost:8025`
* **Access to databases:**
    * PostgreSQL users: `localhost:5433`
    * PostgreSQL admins: `localhost:5434`
    * MongoDB: `localhost:27017`

## Troubleshooting

If you have problems starting the environment, check the container logs using the command `docker logs <container_name>`.
