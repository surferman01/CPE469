# Day 6, 1/28/2025

## Consensus / Election Algorithm

- Used when all network nodes are unaware of a master or if the previous master has gone offline
- How to determine what will be the new master, and spread that information to the entire node network

- Uses
  - Choosing Master
  - Agreeing on value
  - Mutual Exclusion
  - Log Replication

- Ring + Bully algorithms

### Paxos

- a little hard to understand since it was written by some dude
- Requirements
  - Only one request at a time
  - Only one value chosen
  - Master is only known once election function returns

- Proposers
  - Proposes which nodes to become master
- Acceptors
  - Those who can participate in the election
- Learners
  - Entire population; all the nodes

- Messages can take arbitrarily long, but cannot get corrupted

Sequence

- Phase 1
  - Proposer select a proposal numbered n and sends prepare request with number n to a majority of acceptors
  - If an acceptor receives prepare request with a number n greater than any prepare request it has already responded to, then it responds with a promise not to accept any more proposals numbered less than n

- Phase 2
  - If proposer receives a respose to its prepare request (numbered n) from a majority of acceptors, then it sends accept request to each of those acceptors for a proposal numbered n with a value v, where v is the value of the highest numbered proposal among the responses, or any value if the responses reported no proposals
  - if an acceptor receives an accept request for a proposal numbered n, it accepts the proposal unless it has already responded to a prepare request having a number greater than n
