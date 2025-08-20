# Token Swap Tracker

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

## Local Dev

### Install docker and tilt

#### macOS

- Install [Docker for Mac](https://docs.docker.com/desktop/mac/install/)
- Install tilt with:

  ```bash
  curl -fsSL https://raw.githubusercontent.com/tilt-dev/tilt/master/scripts/install.sh | bash
  ```

#### Linux

- Install [Docker](https://docs.docker.com/get-docker/)
- Setup Docker as a [non-root user](https://docs.docker.com/engine/install/linux-postinstall/)
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
