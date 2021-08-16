#!/bin/bash

RECOMMENDED_LAUNCH_FILE=recommended-launch.json
LAUNCH_FILE=launch.json
CURRENT_TIME=$(date "+%Y%m%d-%H%M%S")
LAUNCH_FILE_BACKUP=launch.backup.$CURRENT_TIME.json

if test -f "$LAUNCH_FILE"; then
    cp $LAUNCH_FILE $LAUNCH_FILE_BACKUP
    rm $LAUNCH_FILE
fi
cp $RECOMMENDED_LAUNCH_FILE $LAUNCH_FILE
