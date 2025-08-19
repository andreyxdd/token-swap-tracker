# init

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
