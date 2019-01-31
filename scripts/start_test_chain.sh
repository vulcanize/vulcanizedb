#!/bin/bash

MNEMONIC_PHRASE="whisper ordinary mystery awesome wood fox auction february blind volcano spare soft"
PORT=7545
DATABASE_PATH=test_data/test_chain/
echo Starting ganache chain on port $PORT...

ganache-cli --port $PORT \
            --db $DATABASE_PATH \
            2>&1 > ganache-output.log &
