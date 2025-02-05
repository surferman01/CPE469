# Paxos Notes (2/5/2025)

## Intro

- Paxos = consensus algorithm

## Consensus Algorithm

- how do we support a system that shares/updates information over time?
- need to meet 3 criteria
  1. only a proposed value should be used to update data
  2. only a single value should be elected to update data
  3. process never learns that data is updated unless it actually is for certain

## Characters

1. Proposer
2. Acceptor
3. Learner

- assume communication exists between all nodes
- assume any node may fail
- assume arbitrary code execution speed
- assume messages may take arbitrary time + duplicated + lost, but not corrupted

## Reasoning Process

- Want to send data to nodes
- send data, and if its accepted by majority, then assume it is sent properly
  - this requires that each node can only accept ONE value
  - this is problematic since there may be a tie between multiple messages/data being emitted
- to avoid this problem, give each proposal an integer identifier
  - This can be achieved by having some global integer counter that ensures each proposal gets a unique id or through some other means (implementation based)
- Use the identifier for a couple reasons...

## Proposer Algorithm

Prepare Request

When a node has data to send (proposer):

- create a new proposal with proposal # n
- send request to some subset of nodes (acceptors) asking for:
  - a promise to not accept any proposal with a lower n value
  - reply with the previous proposal (if any) that it has replied to

if the proposer gets a majority of reponses from acceptors:

- issue a proposal with number n and value v to some set (doesn't need to be the same set) of acceptors
- this is called an Accept Request
  - v = v of highest n received from reply to prepare request
    - since each acceptor will send back its highest n prepare request, take the highest n request and use its v as the new v

## Acceptor

When a node receives a Prepare Request:

Can recieve:

- Prepare requests
  - Always ok to respond to
- Accept requests
  - only ok to respond if it hasn't promised not to
  - aka, only if it hasn't responded to a prepare request with higher n than the incoming accept request

## Optimization

- acceptors only need to store highest n proposal its accepted + the highest n prepare request its responded to

## Distributing Data

- many ways to accomplish this
- designate 'distinguished learner'
  - make all acceptances from acceptors route through this node
  - this node then distributes data to all learners

- designate 'distinguished proposer' if no proposals are going through due to too many proposals
  - make proposal requests route through this node so acceptors only speak to this node (or a few of them)
  - could create groups of proposers where each group has a distinguished proposer that actually sends prepare requests out

## Implementation

- primary requirement = ensure no proposal has same n
  - could use disjoint sets for each proposer (like evens / odds) for 2 proposers

## My response to someone elses post

Hey,

Nice summary. I think you are missing a crucial aspect of the Paxos algorithm that is important for maintaining up-to-date information and data though: when a proposer receives the responses from the acceptors, it updates its proposal v value based on what the acceptors send. The acceptor that sends the highest n proposal to the proposer will have the v from that proposal hijacked and replace the v that the original proposal has. Thus, in the time it takes for all the acceptors to reply, if the data changes, then we will be able to get the most updated v data value. Let me know if I am mistaken or if my message is unclear.