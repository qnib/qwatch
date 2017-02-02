# qwatch [![Build Status](http://wins.ddns.net:8000/api/badges/qnib/qwatch/status.svg)](http://wins.ddns.net:8000/qnib/qwatch) [![Coverage](http://wins.ddns.net:8008/badges/qnib/qwatch/coverage.svg)](http://wins.ddns.net:8008/qnib/qwatch)

A golang ETL to collect, filter and output logs/events.

## ROADMAP

### 0.7.x (Inventory)

The current development is aiming to implement a Neo4j-backed inventory ProofOfConcept.

It will use the information fetched from the inputs of `qwatch` and derived inventory information out of it.

- [x] **0.7.0.x**  implement basic Neo4j backend 	

- [x] **0.7.1.x** implement deriving inventory from `docker-events` input
    - [x] Images
    - [x] Containers
    - [ ] network configuration
- [ ] **0.7.2.x** implement deriving inventory from `docker-logs` input
    - [ ] which processes are running inside
    - [ ] how are the files doing
- [ ] **0.7.3.x** implement `docker-engine` log input
- [ ] **0.7.4.x** fetch imput from `Docker SWARM`
