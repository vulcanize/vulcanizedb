# Docker

## headerSync
Dockerfile for running headerSync in a container

### Build
From project root directory:
```
docker build -f dockerfiles/header_sync/Dockerfile . -t header_sync:latest
```

### Run
```
docker run -e DATABASE_USER="user" -e DATABASE_PASSWORD="pw" -e DATABASE_HOSTNAME="host" -e DATABASE_PORT="port" -e DATABASE_NAME="name" -e STARTING_BLOCK_NUMBER=0 -e CLIENT_IPCPATH="path" -t header_sync:latest
```

## resetHeaderCheckCount
Dockerfile for resetting the `headers.check_count` to zero in the database for the given header, so that the execute command
will transform the associated events. This is useful in case an event log is missing.

### Build
From project root directory:
```
docker build -f dockerfiles/reset_header_check_count/Dockerfile . -t reset_header_check_count:latest
```

### Run
```
docker run -e CLIENT_IPCPATH=ipc_path -e DATABASE_USER=user -e DATABASE_PASSWORD=password -e DATABASE_HOSTNAME=host -e DATABASE_PORT=port -e DATABASE_NAME=name -e HEADER_BLOCK_NUMBER=0 -t reset_header_check_count:latest
```
Notes:
- `HEADER_BLOCK_NUMBER` variable is required
- to allow for the command to connect to a database running on the local host, you'll need to:
    - if on MacOS use `host.docker.internal` as the `DATABASE_HOST`
    - if on Linux include the following flag: `--network=host`
