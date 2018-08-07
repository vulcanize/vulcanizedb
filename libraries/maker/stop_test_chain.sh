#!/bin/bash

echo "Stopping ganache chain on port 7545"
ps -ef | grep ganache | grep -v grep | awk '{print $2}' | xargs kill