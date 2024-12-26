# Admin Service Documentation

## Description

This service is responsible for managing administrator accounts and their authentication within the system. It provides a REST API for registering, resetting passwords, and updating administrator data, as well as gRPC for validating administrators.

## Authentication Methods

The service uses JWT (JSON Web Token) to authenticate administrators. Tokens are generated after successful login and must be included in each request to protected resources. In addition, the service uses a mechanism for rotating RSA keys using HashiCorp Vault, which increases security.

## REST API

### Registration

* **Endpoint:** `/api/admin/register`
* **Method:** POST
* **Input data:** JSON with email address, password, and administrator name.
* **Output data:** JSON with the ID of the newly created administrator.

### Password reset

* **Endpoint:** `api/admin/reset-password`
* **Method:** POST
* **Input data:** JSON with the administrator's email address.
* **Output data:** JSON with information about sending a password reset link to the given email address.

### Administrator data update

* **Endpoint:** `/api/admin/update/{admin_id}`
* **Method:** PUT
* **Input data:** JSON with the new email address and administrator name.
* **Output data:** JSON with information about the successful update.

### Administrator removal

* **Endpoint:** `/api/admin/delete/{admin_id}`
* **Method:** DELETE
* **Input data:** Administrator ID in the request path.
* **Output data:** JSON with information about successful deletion.

## gRPC

### Administrator validation

* **Method:** ValidateAdmin
* **Input data:** ValidateAdminRequest with the administrator's email address and password.
* **Output data:** ValidateAdminResponse with information about the validity of the data and the administrator's ID.

## Additional information

* The service stores administrator data in a PostgreSQL database.
* Administrator passwords are hashed using bcrypt.
* The service uses middleware to authorize requests and verify JWT tokens.
* Public keys for JWT token verification are stored in HashiCorp Vault and periodically rotated.
