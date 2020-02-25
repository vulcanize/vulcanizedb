#! /usr/bin/env bash

set -e

function message() {
    echo
    echo -----------------------------------
    echo "$@"
    echo -----------------------------------
    echo
}

ENVIRONMENT=$1
if [ "$ENVIRONMENT" == "prod" ]; then
TAG=latest
elif [ "$ENVIRONMENT" == "staging" ]; then
TAG=staging
else
   message UNKNOWN ENVIRONMENT
fi

if [ -z "$ENVIRONMENT" ]; then
    echo 'You must specifiy an envionrment (bash deploy.sh <ENVIRONMENT>).'
    echo 'Allowed values are "staging" or "prod"'
    exit 1
fi

message BUILDING HEADER-SYNC
docker build -f dockerfiles/header_sync/Dockerfile . -t makerdao/vdb-headersync:$TAG

message BUILDING EXTRACT-DIFFS
docker build -f dockerfiles/extract_diffs/Dockerfile . -t makerdao/vdb-extract-diffs:$TAG

message BUILDING RESET-HEADER-CHECK
docker build -f dockerfiles/reset_header_check_count/Dockerfile . -t makerdao/vdb-reset-header-check:$TAG

message LOGGING INTO DOCKERHUB
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USER" --password-stdin

message PUSHING HEADER-SYNC
docker push makerdao/vdb-headersync:$TAG

message PUSHING EXTRACT-DIFFS
docker push makerdao/vdb-extract-diffs:$TAG

message PUSHING RESET-HEADER-CHECK
docker push makerdao/vdb-reset-header-check:$TAG

# message DEPLOYING SERVICE
