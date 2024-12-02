# Secure JWT Microservices Deployment with Vault and Google Cloud KMS

## Overview

This project implements a secure infrastructure for managing secrets in a microservices architecture, utilizing **Vault**, **Google Cloud KMS**, and **Terraform**. The system is designed to securely handle sensitive data such as unseal keys, policies, and cryptographic keys, which are dynamically delivered to microservices within a Kubernetes (GKE) cluster. Please note that this setup is intended for development purposes and is not production-ready, as it lacks production-grade SSL certificates.

## Architecture & Workflow

### Key Management:
- **Unseal Keys** are securely stored in **Google Cloud KMS** and are delivered to Vault running in the GKE cluster via a **CSI Driver**.
- Other sensitive secrets, such as policies and cryptographic keys (private/public keys), are also stored within Vault.

### Vault Configuration:
- Vault is configured to dynamically supply secrets to microservices on demand, ensuring that only the necessary data is provided securely at runtime.
- The configuration includes integration with **Google Cloud KMS** for key management and encryption.

### Infrastructure Setup:
- **Terraform** is used for provisioning the infrastructure, which includes the setup of Vault and its integration with **Google Cloud KMS**.
- The infrastructure is deployed within a GKE cluster, which is running on **two nodes** to ensure basic availability.

### CI/CD Pipeline:
- The deployment pipeline is managed via **Jenkins**, with encrypted credentials securely stored in **Jenkins Credentials**.
- The pipeline automates the deployment process, ensuring consistency and security through Terraform-managed infrastructure provisioning.

## Microservices Overview

The application is composed of the following microservices, each with a specific responsibility:

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

Each microservice leverages **JWT tokens** for secure communication and authentication.

## Security Design

The security model in place uses best practices for managing secrets and ensuring secure communication:

- **Vault** acts as the central secret management system, enforcing policies, delivering secrets dynamically, and enabling full audit capabilities.
- **Google Cloud KMS** is used to encrypt and securely deliver the unseal keys to Vault, ensuring a high level of security for sensitive data.
- **Jenkins** automates the CI/CD pipeline, ensuring encrypted credentials are used throughout the deployment process.
- **Terraform** ensures infrastructure provisioning is automated and consistent, removing human error in the setup process.

## Advantages

- **Centralized Secret Management**: Securely store and dynamically deliver secrets as needed.
- **High Security**: With Vault and Google Cloud KMS, this setup provides the highest level of security for key management.
- **Scalability**: Built on Kubernetes and Vault, the infrastructure is easily scalable to meet the needs of large, distributed systems.

## Deployment Environment

- The application is deployed within a **Google Kubernetes Engine (GKE)** cluster, spanning **two nodes**. This configuration provides basic availability but is not suitable for production environments.
- The system is not production-ready as it lacks SSL certificates for secure communications, which is a crucial element for deployment in production.

## Local Setup

To run the project locally:

1. Clone the repository:

```bash
git clone https://github.com/nesistor/enigneering_work
cd engineering_work
make_up
