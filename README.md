# Secure JWT Microservices Deployment with Vault and Google Cloud KMS

## Overview

This project implements a secure infrastructure for managing secrets in a microservices architecture using **Vault**, **Google Cloud KMS**, and **Terraform**. It ensures that sensitive data like unseal keys, policies, and cryptographic keys are stored securely and dynamically delivered to microservices in a Kubernetes (GKE) cluster.

## Architecture & Workflow

1. **Key Management:**
   - **Unseal Keys** are stored in **Google Cloud KMS** and delivered to Vault in the GKE cluster via **CSI Driver**.
   - Additional secrets (e.g., policies, private/public keys) are also stored in Vault.

2. **Vault Configuration:**
   - Vault dynamically supplies secrets to the microservices as needed, ensuring secure and on-demand access.

3. **CI/CD Pipeline:**
   - Deployment is automated through a **Jenkins pipeline** with encrypted credentials stored in **Jenkins Credentials**.
   - **Terraform** is used for infrastructure provisioning, including Vault setup and integration with Google Cloud KMS.

## Microservices

The application consists of the following microservices, each with a specific responsibility:

| Service Name      | Description                              |
|-------------------|------------------------------------------|
| **auth-service**  | Authentication service.                  |
| **user-service**  | User management and profile handling.    |
| **admin-service** | Admin functionalities and user control.  |
| **mail-service**  | Email notifications and message handling.|
| **logger-service**| Logging, monitoring, and auditing.       |
| **flyway**        | Database migrations and version control. |
| **vault**         | Key management in the cluster.           |
| **postgresql**    | Relational database for persistent data. |
| **mongodb**       | NoSQL database for unstructured data.    |
| **mailhog**       | Email capturing and testing tool.        |

Each service uses **JWT tokens** for secure communication and authentication.

## Security Design

- **Vault** serves as a centralized secret management system, ensuring dynamic secret delivery, auditing, and policy enforcement.
- **Google Cloud KMS** ensures unseal keys are encrypted and securely delivered to Vault.
- **Jenkins** and **Terraform** provide secure, automated deployment and infrastructure setup.

## Advantages

- **Centralized Secret Management:** Secure storage, dynamic delivery, and auditing of secrets.
- **High Security:** Google Cloud KMS and Vault’s robust features ensure the highest security standards.
- **Scalability:** Easily scalable with Kubernetes and Vault for large, distributed systems.

## Local Setup

Clone the repository:

```bash
git clone https://github.com/nesistor/enigneering_work

Go to the project directory

```bash
  cd enigneering_work/project
```

Start the server

```bash
  make up_build
```

## 🚀 About Me
I'm a full stack developer focusing on Golang and Flutter.

https://www.linkedin.com/in/karol-malicki/
