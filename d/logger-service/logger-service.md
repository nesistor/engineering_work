# Microservice Logging Documentation

This microservice is designed to collect and store logs from other microservices within an application. It offers multiple integration methods to facilitate easy logging from various platforms and programming languages.

## Features

*   **Log Recording:** The microservice allows recording logs from any source.
*   **Log Storage:** Logs are stored in a MongoDB database, ensuring scalability and reliability.
*   **Multiple Integration Methods:** The microservice offers integration via HTTP, gRPC, and RPC, enabling logging from different platforms and programming languages.
*   **Easy Configuration:** The microservice is easy to configure and deploy.

## Integration

The microservice offers three integration methods:

### HTTP

*   Send a POST request to the `/log` endpoint with a JSON payload containing the fields `name` (log name) and `data` (log message).

    *   Example:

    ```json
    {
        "name": "my-service",
        "data": "Log from my service"
    }
    ```

### gRPC

*   Use the gRPC protocol defined in the `logs.proto` file.
*   Call the `WriteLog` method with a `LogRequest` object containing a `Log` object with the fields `name` and `data`.

### RPC

*   Use the RPC protocol.
*   Call the `LogInfo` method with an `RPCPayload` object containing the fields `name` and `data`.

## Architecture

The microservice is built in Go and uses the following libraries:

*   `net/http`: For handling HTTP requests.
*   `google.golang.org/grpc`: For handling gRPC.
*   `net/rpc`: For handling RPC.
*   `go.mongodb.org/mongo-driver`: For interacting with the MongoDB database.
*   `github.com/go-chi/chi/v5`: For HTTP routing.

## Deployment

The microservice can be deployed on any platform that supports Go.

## Configuration

The microservice can be configured via environment variables:

*   `webPort`: HTTP server port (default 80).
*   `gRpcPort`: gRPC server port (default 50003).
*   `rpcPort`: RPC server port (default 5001).
*   `mongoURL`: MongoDB database URL (default `mongodb://mongo:27017`).

## Monitoring

The microservice does not have a built-in monitoring system. It is recommended to use external tools to monitor the health of the microservice and the database.

## Scalability

The microservice is scalable due to the use of MongoDB. It can be deployed in a cluster and load balanced across multiple instances.

## Remarks

*   The microservice does not support authentication or authorization.
*   The microservice does not have a log rotation mechanism.
*   The microservice does not offer a user interface for viewing logs.

## Further Development

In the future, the following features are planned to be added:

*   Authentication and authorization.
*   Log rotation.
*   User interface for viewing logs.
*   Support for additional log formats.

## Contact

If you have any questions or problems, please contact the system administrator.