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
