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

# build images
message BUILDING HEADER-SYNC
docker build -f dockerfiles/header_sync/Dockerfile . -t makerdao/vdb-headersync:$TAG

message BUILDING RESET-HEADER-CHECK
docker build -f dockerfiles/reset_header_check_count/Dockerfile . -t makerdao/vdb-reset-header-check:$TAG

message LOGGING INTO DOCKERHUB
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USER" --password-stdin

# publish
message PUSHING HEADER-SYNC
docker push makerdao/vdb-headersync:$TAG

message PUSHING RESET-HEADER-CHECK
docker push makerdao/vdb-reset-header-check:$TAG

# service deploy
if [ "$ENVIRONMENT" == "prod" ]; then
  message DEPLOYING HEADER-SYNC
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-header-sync-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$PROD_REGION.amazonaws.com --region $PROD_REGION

elif [ "$ENVIRONMENT" == "staging" ]; then
  message DEPLOYING HEADER-SYNC
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-header-sync-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$STAGING_REGION.amazonaws.com --region $STAGING_REGION
else
   message UNKNOWN ENVIRONMENT
fi
