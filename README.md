# A Golang-based REST API service

## Design
The application was implemented according to the provided specifications. The main components are:
- Programming language: Go, selected for its strong typing and concurrency model
- Storage: SQLite, selected for its simplicity and performance
- Web server: Gin, selected for its lightweight HTTP routing and middleware support
- Logging: Zerolog, selected for its low overhead and efficient structured logging (useful for log analisys tools like Datadog)

## Endpoints
- GET `/objects/{bucket}/{objectID}`: retrieve an object by bucket ID and object ID
- PUT `/objects/{bucket}/{objectID}`: create an object in the specified bucket using the provided object ID
- DELETE `/objects/{bucket}/{objectID}`: delete an object by bucket ID and object ID

## Environment variables
- `LISTENING_ADDRESS`: IP address or hostname on which the web server listens, defaults to `0.0.0.0`
- `HTTP_PORT`: TCP port on which the web server listens, defaults to `8080`

## Run tests
To run the integration test suite, execute `make test`.
> NOTE: due to the limited scope of the task, only integration tests were implemented.

## Deploy
To deploy locally, run `make all` or `make all-docker`.
To deploy to production, apply the Kubernetes manifests in the `deploy` directory.

## Production Notes
- Migrate to an external PostgreSQL instance as the primary datastore
- Add a cache layer (e.g. Redis)
- Add a `/metrics` endpoint
- Add Swagger and OpenAPI documentation
- Move to a fully RESTful resource structure `/buckets/{bucketID}/objects/{objectID}`:
    - GET object on the `/buckets/{bucketID}/objects/{objectID}` endpoint
    - POST new objects on the `/buckets/{bucketID}/objects` endpoint (the object ID must be assigned by the backend)
    - PUT existing objects (override) on the `/buckets/{bucketID}/objects/{objectID}` (or just use PATCH)
    - DELETE object on the `/buckets/{bucketID}/objects/{objectID}` endpoint
