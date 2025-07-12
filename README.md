# Device Management Service

This project demonstrates a horizontally scalable backend service in Go, using [Temporal](https://temporal.io/) for robust workflow orchestration. It is designed to support mass device configuration tasks, suitable for deployment on Kubernetes, and provides rapid local development and testing cycles.

## Features

- **Go Backend:** Fast, statically compiled, and easy to run locally or in containers.
- **Temporal Workflows:** Reliable orchestration for complex automation like mass device provisioning/configuration.
- **Automated Testing:** Unit and workflow tests using Go's `testing` and Temporal's test suite.
- **HTTP API:** REST endpoints to trigger and monitor workflows.
- **Kubernetes Ready:** Easily containerizable and deployable.

## Project Structure

```
dm-backend/
├── cmd/
│   └── server/                  # Application entry point
├── internal/
│   ├── api/                     # HTTP API handlers & server setup
│   ├── activities/              # Temporal workflow activities
│   ├── workflow/                # Workflow orchestration logic
├── pkg/
│   └── models/                  # Shared domain models
├── go.mod
└── README.md
```

## Running Locally

1. **Start a Temporal server** (see [Temporal Docker Compose](https://docs.temporal.io/v1.0/docs/server/docker-compose/)).
   
2. **Temporal Server Configuration**

By default, the service connects to Temporal at `localhost:7233` for development.

To override, set the `TEMPORAL_HOSTPORT` environment variable:

```bash
export TEMPORAL_HOSTPORT="temporal.example.com:7233"
go run cmd/server/main.go
```

If `TEMPORAL_HOSTPORT` is not set, the service will use `localhost:7233`.

3. **Run the worker and API server:**
   ```bash
   go run cmd/server/main.go
   ```
4. **Test API endpoints:**
   - Start a workflow:
     ```bash
     curl -X POST http://localhost:8080/api/config/start \
       -H "Content-Type: application/json" \
       -d '{"devices":[{"ID":"dev1","Name":"Device 1"}],"configParams":{"mode":"fast"}}'
     ```
   - Check workflow status:
     ```bash
     curl "http://localhost:8080/api/config/status?workflowID=<workflowID>&runID=<runID>"
     ```

## Testing

Run all tests:
```bash
go test ./...
```

## Example Workflow: Mass Device Configuration

- Accepts a list of devices and configuration parameters.
- Configures each device in parallel using Temporal activities.
- Tracks success/failure for each device.

## Extending

- Add more device activities to `internal/activities/`.
- Add new workflows to `internal/workflow/`.
- Add more API endpoints to `internal/api/`.
- Integrate with external systems or device APIs as needed.

## Deployment

- Add a `Dockerfile` for containerization.
- Use Helm charts or Kubernetes manifests to deploy on K8s.

---

**Questions or improvements? Open an issue or PR!**