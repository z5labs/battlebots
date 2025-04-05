---
title: Proof Of Concept
type: docs
no_list: true
---

## Goals

- Be able to perform a basic action
- Communicate state events to another game server
- Collect telemetry on entire end-to-end flow of actions and state events
- At the end, telemetry data should be used to evaluate effectiveness of this POC

## Tech Stack

- [gRPC](https://grpc.io/) - for game server API
- [ValKey Streams](https://valkey.io/topics/streams-intro/) - for distributed pub/sub between the game servers
- [OTel Collector](https://opentelemetry.io/docs/collector/) - for ingesting telemetry signals from everything
- [InfluxDB](https://www.influxdata.com/) - for storing and indexing telemetry signals from collector
- [Grafana](https://grafana.com/) - for visualizing telemetry signals

## Design

```mermaid
architecture-beta
    group terminal1(cloud)[Terminal 1]
    service gameserver1(server)[Game Server] in terminal1

    group terminal2(cloud)[Terminal 2]
    service gameserver2(server)[Game Server] in terminal2

    group terminal3(cloud)[Terminal 3]
    service cli1(server)[CLI] in terminal3

    group terminal4(cloud)[Terminal 4]
    service cli2(server)[CLI] in terminal4

    group terminal5(cloud)[Terminal 5]
    service valkey(server)[ValKey] in terminal5

    cli1{group}:T <--> B:gameserver1{group}
    gameserver1{group}:R <--> L:valkey{group}
    cli2{group}:T <--> B:gameserver2{group}
    gameserver2{group}:L <--> R:valkey{group}
```

### Observability

```mermaid
architecture-beta
    service cli(server)[CLI]
    service gameserver(server)[Game Server]
    junction signals

    group observability(cloud)[Observability]
    service collector(server)[OTel Collector] in observability
    service influxdb(server)[InfluxDB] in observability
    service grafana(server)[Grafana] in observability

    collector:T --> B:influxdb
    grafana:R --> L:influxdb

    cli:T -- B:signals
    gameserver:B -- T:signals
    signals:L --> R:collector
```

## Steps

### Step 0

Define game server interface with support for:

- One action, `Move`
- One state event, `Position`

### Step 1

Implement CLI tool for manually performing `Move` action
and receiving `Position` updates.

### Step 2

Implement a game server with support for one action, `Move`, and
one state event, `Position`.

### Step 3

Run ValKey in local container runtime

### Step 4

Integrate game server with ValKey

### Step 5

Run OTel collector, InfluxDB, Grafana locally

### Step 6

Integrate CLI and game server with OTel collector

### Step 7

In multiple local terminals, run 2 game servers and use 2 CLIs to perform `Move`s

## Results

*TBD: once the POC is completed screenshots and a write up should be put here*