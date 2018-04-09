#!/bin/bash

LOGFILE="$/opt/gocode/src/github.com/MDLlife/teller/scripts/eth-ssh.log"
NOW="$(date +%d/%m/%Y' - '%H:%M)"


REMOTE_USER=teller
REMOTE_HOST=_YOUR_TARGET_SERVER_HERE_

SSH_REMOTE_PORT=22
SSH_LOCAL_PORT=22

TUNNEL_REMOTE_PORT=8545
TUNNEL_LOCAL_PORT=8545

TUNNEL_SPECIFIC_REMOTE_PORT=22

createTunnel() {
    /usr/bin/ssh -f -N  -L$TUNNEL_LOCAL_PORT:0.0.0.0:$TUNNEL_REMOTE_PORT $REMOTE_USER@$REMOTE_HOST -p $TUNNEL_SPECIFIC_REMOTE_PORT
    if [[ $? -eq 0 ]]; then
        echo [$NOW] Tunnel to $REMOTE_HOST created successfully » $LOGFILE
    else
        echo [$NOW] An error occurred creating a tunnel to $REMOTE_HOST RC was $? » $LOGFILE
    fi
    }

## Run the 'ls' command remotely.  If it returns non-zero, then create a new connection
/usr/bin/ssh -p $SSH_LOCAL_PORT $REMOTE_USER@REMOTE_HOST ls >/dev/null 2>&1
if [[ $? -ne 0 ]]; then
    echo [$NOW] Creating new tunnel connection » $LOGFILE
    createTunnel
fi

# netstat -lnpt | awk '$4 ~ /:8545$/ {sub(/\/.*/, "", $7); print $7}'  |awk '{print $1}'| xargs kill | true;

# crontab -e
# */5 * * * *  /opt/gocode/src/github.com/MDLlife/teller/scripts/teller-remote.sh