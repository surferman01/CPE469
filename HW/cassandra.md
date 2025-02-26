# Cassandra Notes (2/24/2025)

Decentralized structured storage system

## Overview

- distributed storage system that promotes high availability of data
- high writethrough w/o sacrificing read efficiency
- used as backend storage system for multiple services in Facebook

## Problem

- Facebook has immense amt of data w strict data integrity requirements
- support continuous growth (scalability) and be fail tolerant

- Inbox Search Problem
  - allows searching through facebook inbox
  - this requires very high write throughput (billions of writes p day)
  - being able to replicate and distribute data with low latency for users

## Data Model

- multi-dimensional map indexed by key
- associated value = highly structured object
- row key = string identifier for an entire row of multi-dimensional columns
  - each operation performed on a row key is atomic (per replica)
- column families
  - columns grouped together into sets
  - super:
    - a column family inside another column family
  - applications can specify sort order of column families

## API (I/O)

insert(table, key, rowMutation)
get(table, key, columnName)
delete(table, key, columnName)

- columnName is very flexible and can be:
  - specific column within a column family
  - column family
  - super column family
  - column within a super column

## Architecture

- abides by and requires typical distributed system requirements
- notable changes:
  - partition handling, replication, membership, failure tolerance/handling, scaling

- quick overview:

for writes:
sys routes request to multiple replicas and waits for ack

for reads:
either route to nearest (geographically) server for data, or route to all replicas and wait for many responses (depending on data requirements)

### Partitioning

Performed and managed by a consistent hashing scheme

- uses order preserving hash function

Impelments the typical consistent hashing algorithm

- for load balancing, it periodically assesses the load balance and re-distributes nodes to make it even

### Replication

- Replicate data on n (chosen) replicas
- each key k is assigned a coordinator who ensures replication of data on the range of keys it monitors
- each row should be replicated across other data centers -> more fail tolerance

### Membership

- based on Scuttlebutt
- anti-entropy Gossip based protocol
  - good because of efficient CPU utilization
- gossip is also used in Cassandra to "disseminate other system related control states"

Failure Detection

- uses: Accrual Failure Detection
  - instead of saying node is dead/alive, it gives 'suspicion' value to other nodes
  - user defines functionality of this suspicion value
  - good for changing network conditions and highly accurate

### Bootstrapping

- on startup, randomly choose a position on the ring
- on startup error, manager can manually configure the node

### Scaling

- its easy to add nodes and move around nodes via Cassandra interface / cli

### Local Persistence

- recall, most of the data/information lives locally on the nodes
- this requires reads/writes to disk
- use a commit log to track adds/removes/updates
- write to main memory only when some threshold of temp data is achieved
- sequential writes with an index based on row key for quick lookup

### Implementation Details

- Cassandra process components on single machine:
- partitioning model
- cluster membership
- failure detection module
- storage engine model

- membership module on top of network layer w non-blocking I/O (so its lightweight)
- control messages are passed thru UDP
- application related messages passed thru TCP (replication/requests)

Upon request to node, enter following states:

1. identify node with pertinent information/data for key
2. route request to node and await response
3. if they dont arrive w/in timeout, fail request and return to client
4. determine latest response based on timestamp
5. schedule repair of data for replica(s) w/o data

- the commit log is basically just a stack w fixed length, so older data just gets overwritten
- if data has been written to disk, purge commit log
- for writing to disk, use fast sync or normal
  - fast sync will buffer the data and allow regular operation of node simultaneously
    - if the node crashes while data in buffer, then data is lost
  - normal will sequentially write data to disk
