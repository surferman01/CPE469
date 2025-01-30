# Day 7 (1/30/2025)

Leader Election Algorithms

## Paxos

### Steps

1. A proposer select a proposal numbered n and sends prepare request with number n to majority of acceptors
2. If an acceptor receives prepare request with number n greater than any request that it has already responded to, then it responds with a promise not to accept any more proposals numbered less than n or with the highest numbered proposal thats accepted
3. If proposer receives a response to its prepare request (n)

## Raft

1. Give each candidate a random timeout timer
2. on timeout, send to all nodes that you will be leader
3. If you get the majority of responses, then you become leader
    - each node can only accept one election
4. if there are multiple timers at same time and no leader is selected, just re-run the system

- Each new leader that gets elected gets an election number (incremented by 1 for each new leader)
  - If the network gets partitioned, then when they recombine, the highest election number will take priority

for Log matching:

- from leader, send most recent data along with new data
- compare recent data to see if its the same
  - if yes, then write new data
  - if no, then second most recent data and compare the two to see if same.
    - if yes, then write recent and new data
    - if no, repeat

## Consensus Algorithm

### 3 uses

1. Leader election
2. (Log) Replication
3. Safety
