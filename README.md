# qwatch [![Build Status](http://wins.ddns.net:8000/api/badges/qnib/qwatch/status.svg)](http://wins.ddns.net:8000/qnib/qwatch) [![Coverage](http://wins.ddns.net:8008/badges/qnib/qwatch/coverage.svg)](http://wins.ddns.net:8008/qnib/qwatch)

A golang ETL to collect, filter and output logs/events.

## ROADMAP

### 0.7.x (Inventory)

The current development is aiming to implement a Neo4j-backed inventory ProofOfConcept.

It will use the information fetched from the inputs of `qwatch` and derived inventory information out of it.

- **0.7.1.x** implement deriving inventory from `docker-events` input
    - [ ] .1 Images
    - [ ] .2 Containers
    - [ ] .3 network configuration
    - [ ] .4 docker-engines
    - [ ] .5 Docker SWARM
- **0.7.2.x** implement deriving inventory from `docker-logs` input
    - [ ] .1 which processes are running inside
    - [ ] .2 how are the files doing
