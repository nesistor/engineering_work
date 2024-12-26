# User Service: Comprehensive Documentation

This document provides a detailed architectural overview of the user service, a critical component of a larger microservices ecosystem. It is responsible for managing core user lifecycle functionalities, including registration, authentication, authorization, and account maintenance.

## Architectural Overview

The user service is designed as a lightweight, stateless microservice, promoting scalability and maintainability. It leverages the Go programming language for its efficiency and concurrency support, coupled with PostgreSQL as the persistent data store for user information. Security is paramount, with integration with HashiCorp Vault for secure storage and rotation of cryptographic keys used in JWT-based authentication.

## Core Functionalities

* **User Registration**: New users can register with the service, providing essential information such as email, username, and password. Data validation is enforced to ensure data integrity and security.
* **Authentication**: Users can authenticate using their registered credentials. The service employs JWT (JSON Web Tokens) as the authentication mechanism, providing a secure and stateless approach to user verification.
* **Authorization**: The service supports role-based access control (RBAC). Users are assigned roles, and access to resources is granted based on these roles, ensuring only authorized users can perform specific actions.
* **Account Management**: Registered users can manage their accounts, including updating personal information and passwords. The service provides APIs for these operations, allowing users to control their data.
* **Password Reset**: A robust password reset mechanism is implemented, enabling users to securely reset their passwords in case of forgetfulness or compromise. This process utilizes time-limited tokens and email notifications to ensure security.

## Microservice Interactions

The user service actively interacts with other microservices in the ecosystem to provide comprehensive functionality:

* **Logger Service**: All significant events, such as user registrations, updates, and deletions, are logged asynchronously to the logger service. This centralized logging approach facilitates monitoring, auditing, and debugging.
* **Mail Service**: The user service utilizes the mail service to send email notifications, including password reset links and other account-related communications.

## Technology Stack

* **Go**: The service is implemented in Go, chosen for its performance, concurrency features, and strong support for microservice development.
* **PostgreSQL**: PostgreSQL is used as the relational database to store user data. Its robustness, reliability, and ACID properties make it suitable for managing critical user information.
* **HashiCorp Vault**: Vault is integrated to provide a secure and centralized solution for managing cryptographic keys. It ensures secure key storage, rotation, and access control, enhancing the overall security posture.
* **gRPC**: gRPC is used for inter-service communication, providing a high-performance, type-safe, and efficient mechanism for microservices to interact.
* **JWT (JSON Web Tokens)**: JWT is employed for stateless user authentication. It allows secure transmission of user identity information between services.

## Deployment and Configuration

The user service is containerized using Docker, ensuring portability and consistency across different environments. It can be deployed on container orchestration platforms like Kubernetes for scalability and resilience.

Configuration is typically managed through environment variables, allowing for easy customization and integration into various deployment pipelines.

## Security Considerations

Security is a primary concern in the design and implementation of the user service. Key security measures include:

* **Password Hashing**: Passwords are stored as strong, one-way hashes using bcrypt, protecting them even in case of database compromise.
* **Secure Key Management**: Vault is used to securely store and manage cryptographic keys, ensuring they are protected from unauthorized access and regularly rotated.
* **Input Validation**: All user inputs are rigorously validated to prevent injection attacks and ensure data integrity.
* **Rate Limiting**: Rate limiting can be implemented to prevent brute-force attacks and protect the service from abuse.

## Future Enhancements

While the user service provides a solid foundation for user management, several potential enhancements can further improve its capabilities and security:

* **Two-Factor Authentication (2FA)**: Implementing 2FA would add an extra layer of security, requiring users to provide a second form of verification, such as a time-based one-time password (TOTP).
* **Fine-Grained Authorization**: More granular authorization mechanisms, such as attribute-based access control (ABAC), could be explored to provide more fine-grained control over user permissions.
* **User Activity Monitoring**: Monitoring user activity, including login attempts, profile updates, and password resets, can help detect suspicious behavior and enhance security.
* **Integration with Identity Providers**: Integrating with external identity providers, such as Google, Facebook, or Okta, can streamline user authentication and improve user experience.

## Conclusion

The user service provides a robust and secure foundation for managing user lifecycles in a microservices architecture. Its careful design and implementation, leveraging technologies like Go, PostgreSQL, Vault, and gRPC, ensure scalability, maintainability, and security. Continuous improvements and enhancements will further solidify its role as a critical component of the overall system.

# Microservice User-Service Documentation

## General Description

The `user-service` microservice is a key element of the system, responsible for managing users. Its main functionalities include:

* Registration of new users
* Authentication of users
* Management of user data (updating, deleting)
* Password reset

## Architecture

The service is implemented in Go and uses the following technologies:

* PostgreSQL as the database
* gRPC for communication with other services
* JWT for authentication
* HashiCorp Vault for key management
* chi as the HTTP router

## Implementation Details

### User Management

The `User` struct in the `models.go` file defines the user model:

```go
type User struct {
    ID           int64     `json:"id"`
    UserName     string    `json:"username"`
    Email        string    `json:"email"`
    PasswordHash string    `json:"passwordhash"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"udpated_at"`
}

### User Management

The `UserModel` model in the same file contains methods for interacting with the database, such as:

* `GetAllUsers()` - returns a list of all users
* `GetUserByEmail()` - returns the user with the given email address
* `InsertUser()` - adds a new user to the database
* `UpdateUser()` - updates user data
* `DeleteUserByID()` - deletes the user with the given ID

### Authentication

Authentication is done using JWT tokens. The keys for generating and verifying tokens are stored in HashiCorp Vault and managed by `KeyManager` (`keys.go` file).

### gRPC Communication

The service provides the gRPC method `ValidateUser`, which is used to verify the user based on email address and password. The definition of the method is located in the `users.proto` file.

### HTTP API

The service also provides an HTTP API with the following endpoints:

* `/api/login/register` - register a new user
* `/api/login/check-email` - check if a user with the given email address exists
* `/api/login/reset-password` - request a password reset
* `/api/login/delete-user/{user_id}` - delete a user (requires authentication)
* `/api/login/update/{user_id}` - update user data (requires authentication)

### Middleware

The `AuthMiddleware` middleware in the `middleware.go` file is used to authenticate HTTP requests using JWT tokens.