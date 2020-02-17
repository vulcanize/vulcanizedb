# Docker

Note: to allow for the command to connect to a database running on the localhost, you'll need to:
- set `DATABASE_HOSTNAME=host.docker.internal` on MacOS
- include the `--network=host` flag on Linux (with `DATABASE_HOSTNAME=localhost`)

For more information, see https://github.com/docker/for-linux/issues/264

## extractDiffs
Dockerfile for populating storage diffs in db to be transformed by `execute`.

## Build
From project root directory:
```
docker build -f dockerfiles/extract_diffs/Dockerfile . -t extract_diffs:latest
```

### Run
Against statediffing Geth pubsub:
```
docker run -e DATABASE_USER=user -e DATABASE_PASSWORD=password -e DATABASE_HOSTNAME=host -e DATABASE_PORT=port -e DATABASE_NAME=name -e CLIENT_IPCPATH=path -e STORAGEDIFFS_SOURCE=geth -it extract_diffs:latest
```

Against CSV:
```
docker run -e DATABASE_USER=user -e DATABASE_PASSWORD=password -e DATABASE_HOSTNAME=host -e DATABASE_PORT=port -e DATABASE_NAME=name -e CLIENT_IPCPATH=path -e FILESYSTEM_STORAGEDIFFSPATH=/data/<csv_filename> -v <csv_filepath>:/data -it extract_diffs:latest
```


## headerSync
Dockerfile for running `headerSync` in a container

### Build
From project root directory:
```
docker build -f dockerfiles/header_sync/Dockerfile . -t header_sync:latest
```

### Run
```
docker run -e DATABASE_USER=user -e DATABASE_PASSWORD=password -e DATABASE_HOSTNAME=host -e DATABASE_PORT=port -e DATABASE_NAME=name -e STARTING_BLOCK_NUMBER=0 -e CLIENT_IPCPATH=path -it header_sync:latest
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
docker run -e CLIENT_IPCPATH=ipc_path -e DATABASE_USER=user -e DATABASE_PASSWORD=password -e DATABASE_HOSTNAME=host -e DATABASE_PORT=port -e DATABASE_NAME=name -e HEADER_BLOCK_NUMBER=0 -it reset_header_check_count:latest
```
Notes:
- `HEADER_BLOCK_NUMBER` variable is required

