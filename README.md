# qwatch [![Build Status](http://wins.ddns.net:8000/api/badges/qnib/qwatch/status.svg)](http://wins.ddns.net:8000/qnib/qwatch) [![Coverage](http://wins.ddns.net:8008/badges/qnib/qwatch/coverage.svg)](http://wins.ddns.net:8008/qnib/qwatch)

A golang ETL to collect, filter and output logs/events.

More information of the graph backend and how it is modeled, can be found [here](Neo4j.md)

## ROADMAP

### 0.7.x (Inventory)

The current development is aiming to implement a Neo4j-backed inventory ProofOfConcept.

It will use the information fetched from the inputs of `qwatch` and derived inventory information out of it.

- [x] **0.7.0.x**  implement basic Neo4j backend 	

- [x] **0.7.1.x** implement deriving inventory from `docker-events` input
    - [x] Images
    - [x] Containers
    - [x] network configuration
- [ ] **0.7.2.x** create `docker-api` input
    - [x] Fetch and push engine `Info()`
    - [x] fetch and push `swarm.Node`
    - [ ] fetch and push `ContainerList` #22
    - [x] fetch and push `ImageList` #23
    - [ ] fetch and push `NetworkList` #24
    - [ ] fetch and push `VolumeList` #25
    - [ ] fetch and push `ServiceList` #26
    - [ ] fetch and push `TasksList` #27
    - [ ] fetch and push `SecretList` #28
    - [ ] fetch and push `PluginList` #29
    - [ ] fetch and push `CheckpointList` #30

- [ ] **0.7.x** implement deriving inventory from `docker-logs` input
    - [ ] which processes are running inside
    - [ ] how are the files doing
