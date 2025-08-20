# Token Swap Tracker

This project simulates the tracking and aggregation of the data from a swap token service.

- `producer` is a Kafka producer service that emits swap events.
- `consumer` contains the implementation of two services:
  - API to serve aggregated data via REST.
  - Kafka consumer that reads data from topics and broadcasts it via a WebSocket server.

Since the Kafka consumer implementation is separate from the REST API service, it can be deployed in a separate Kubernetes deployment and scaled horizontally, depending on the producer's swap event rate or the number of connected WebSocket clients.

## Opened Questions

### What transport mechanisms should be used by the producer?

The Kafka transport mechanism should be sufficient here to handle the 1000 events per minute emitted by the producer. Further performance optimization can be achieved with Protobufs.

### Where should different types of data be stored?

In-memory storage solutions, such as Redis or Memcached, can be utilized to serve data via the REST API. For historical data, time-series databases like InfluxDB or TimescaleDB (a PostgreSQL extension) are the best choices. PostgreSQL or similar databases can be used to store service metadata and checkpoints.

### How can high availability and zero data loss be ensured?

High availability can be achieved in the following ways:
- The Kafka consumer can be scaled horizontally, depending on the load and number of connected WebSockets.
- REST API services can also be scaled horizontally via Kubernetes, as they are stateless.
- Redis can be scaled horizontally using techniques like Redis Cluster, sharding, or read replicas.

To ensure zero data loss, Kafka consumer groups with offset commits after processing can be utilized. The last processed offset can be stored in PostgreSQL as a checkpoint, so that a restart can be done from this saved state.


## Local Dev

### Install docker and tilt

#### Linux

- Install [Docker](https://docs.docker.com/get-docker/)
- Setup Docker as a [non-root user](https://docs.docker.com/engine/install/linux-postinstall/)
- Install tilt with:

  ```bash
  curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
  ```

#### macOS

- Install [Docker for Mac](https://docs.docker.com/desktop/mac/install/)
- Install tilt with:

  ```bash
  curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
  ```

#### Windows

- Install [Docker for Windows](https://docs.docker.com/desktop/windows/install/)
- Install tilt with powershell script:

  ```PowerShell
  iex ((new-object net.webclient).DownloadString('https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.ps1'))

### Launch services

To launch all the services defined in docker-compose, run:
```bash
tilt up
```
Any code updates results in the hot-reloading of the corresponding containers.

## Possible improvements

- Utilize protobufs to further improve serialization in the consumer service
- Handle duplicated events and its order from the producer
- Handle the status of the transcation i.e., pending swap transcations should not be counted towards aggregated stats
- Utilize another Redis instance to keep track of connected clients for web-socket server. Plus, use distributed lock
- Add customized logger (such as zaplog)
- Implement various middlewares in REST API service for rate limiting, auth, TLS, origing checks, etc.
- Further optimize docker images with dockerignore, etc.
- Extract config from environmental variables (such as port, credentials, etc.) or config files, print and validate them on startup
- In WebSocket server, keep write deadlines short and drop slow clients to protect the server
- Utilize separate Kafka topic for each token
