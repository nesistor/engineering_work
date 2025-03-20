# Mail Service Documentation

This document provides a comprehensive overview of a microservice designed to send emails. It includes details about the service's architecture, API endpoints, and how to use it.

## 1. Introduction

The Mail Service is a RESTful API that allows you to send emails using SMTP. It is built using Go and utilizes the `go-simple-mail/v2` library for SMTP interaction and `go-premailer` for inlining CSS into email bodies.

## 2. API Endpoints

The Mail Service exposes the following API endpoint:

* **POST /send**
    * This endpoint is used to send an email.
    * It accepts a JSON payload with the following structure:

    ```json
    {
      "from": "sender@example.com",
      "to": "recipient@example.com",
      "subject": "Email Subject",
      "message": "Email Body"
    }
    ```

    * It returns a JSON response indicating whether the email was sent successfully.

## 3. Code Overview

The Mail Service code is organized into several files:

* **handlers.go:**
    * Contains the handler function for the `/send` endpoint.
    * It reads the JSON payload, constructs an email message, and uses the `Mailer` component to send the email.

* **routes.go:**
    * Sets up the API routes using the `go-chi/chi/v5` library.
    * It defines middleware for CORS and heartbeat checks.

* **mailer.go:**
    * Contains the `Mail` struct, which holds the SMTP configuration.
    * It also contains the `Message` struct, representing an email message.
    * The `SendSMTPMessage` function handles sending the email using the `go-simple-mail/v2` library.
    * It includes functions for building HTML and plain text email bodies from templates.
    * It uses `go-premailer` to inline CSS into the HTML email body.

* **helpers.go:**
    * Contains helper functions for reading and writing JSON responses.
    * It also includes a function for generating error responses.

* **main.go:**
    * Initializes the application and starts the HTTP server.
    * It reads the SMTP configuration from environment variables.

* **handlers_test.go:**
    * Contains unit tests for the handler functions.

## 4. SMTP Configuration

The Mail Service reads its SMTP configuration from the following environment variables:

* `MAIL_DOMAIN`
* `MAIL_HOST`
* `MAIL_PORT`
* `MAIL_USERNAME`
* `MAIL_PASSWORD`
* `MAIL_ENCRYPTION`
* `FROM_NAME`
* `FROM_ADDRESS`

## 5. Dependencies

The Mail Service relies on the following external libraries:

* `go-chi/chi/v5`: For routing and middleware.
* `go-chi/cors`: For CORS support.
* `go-simple-mail/v2`: For sending emails via SMTP.
* `go-premailer`: For inlining CSS into email bodies.

## 6. Running the Service

To run the Mail Service, you need to:

1. Set the required environment variables for SMTP configuration.
2. Build the application using `go build`.
3. Run the executable.

## 7. Testing

The Mail Service includes unit tests for the handler functions. You can run the tests using `go test`.

## 8. Future Improvements

* Add support for more authentication mechanisms (e.g., OAuth2).
* Implement email queuing for better performance and reliability.
* Provide more detailed logging and error handling.
* Add support for sending emails asynchronously.

## 9. Conclusion

The Mail Service provides a simple and effective way to send emails from your applications. It is easy to configure and use, and it can be easily integrated into your existing infrastructure.