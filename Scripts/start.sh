#!/bin/bash

# Exit on first error, print all commands.
set -ev

# Grab the Car Lease Demo directory.
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )";

# If the keyValStore exits, remove it
if [ -f $DIR/keyValStore ]
then
    rm -r $DIR/keyValStore;
fi

# Tag the latest version of fabric-baseimage
docker pull hyperledger/fabric-baseimage:x86_64-0.1.0
docker tag hyperledger/fabric-baseimage:x86_64-0.1.0 hyperledger/fabric-baseimage:latest

#docker pull hyperledger/fabric-peer:latest
#docker pull hyperledger/fabric-membersrvc:latest

# Clean up old docker containers
#docker-compose -f $DIR/docker-compose.yml kill;
#docker-compose -f $DIR/docker-compose.yml down;
#docker-compose -f $DIR/docker-compose.yml build;
#docker-compose -f $DIR/docker-compose.yml up -d ;


docker-compose -f $DIR/docker-compose-four-peer-ca.yaml kill;
docker-compose -f $DIR/docker-compose-four-peer-ca.yaml down;
docker-compose -f $DIR/docker-compose-four-peer-ca.yaml build;
docker-compose -f $DIR/docker-compose-four-peer-ca.yaml up -d ;

