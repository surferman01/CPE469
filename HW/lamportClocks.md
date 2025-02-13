# Lamport Clock Notes

## Problem

Syncing clocks is difficult given clock skew/drift

- How do you know the order of events objectively?

## Partial Ordering

- lets define "happened before" without using physical clocks (->)

3 requirements:

1. if a and b are in the same process and a comes before b, then a -> b
2. if a is sending a message, and b is the receipt of the same message from another process, then a -> b
3. if a -> b and b -> c, then a -> c

- 2 distinct events a and b are concurrent if: a !-> b and b !-> a
- use space-time diagrams to understand ordering
- events are concurrent, even if their physical timing is different, if they do not know exactly what a process is executing
  - ex) in physical time: if a happens before b, but the thread running b last heard from the thread running a before a, then we just know that b and a are concurrent and not exactly when each happened

## Logical Clocks

Define a clock for each event (C)
clock condition:

- if a -> b, then C(a) < C(b)

some implications here:

- if a and b are events in same process where a comes before b, then C(a) < C(b)
- if a is sending a message to another process and b is the receiving of it, then C(a) < C(b)

## Idea

- to sync clocks we need to meet these conditions:

1. each process increments C between two consecutive events
2. if event a is sending a message, then give it a timestamp (C(a))
    - when the message is recieved, if the messages timestamp is greater than the receiver, then increment the reciever clock to be just greater than the message timestamp

Total ordering (=>):

- if a is an event in a process and b in an event in another process, then a => b if:
  1. C(a) < C(b)
  or
  2. C(a) = C(b) and process of a < process of b

## Shared resources

Write algorithm for granting resources to process that satisfies following:

1. process granted resource must release it before its issued to another process
2. requests for resource must be handled in order they occur
3. if every process which is granted resource eventually releases it, then every request is eventually granted

### Solution

- each process has its own request queue

5 rules

1. to request resource, send request to access resources to ALL other processes with timestamp m
2. when a process recieves such a request, put it in its own request queue and send its current timestamp (as ack)
3. to release resource, original process removes any of its own requests for the resource in its request queue and sends release request w timestamp
4. when process recieves release request, remove the entry from its request queue
5. original process gets the resource if the following 2 conditions are true:
  a) your resource request in the requests queue is the lowest timestamp
  b) the original process has heard an ack from each other process timestamped later than the resource request timestamp

## Strong clock condition

if we have some other communication methods and cannot confirm C(a) < C(b),
impose the following:

- if a --> b, then C(a) < C(b)

## My Response to someone

Hey,

I appreciate your concise summary of the concept of logical clocks presented by Leslie Lamport. You mention that a receiver updates its own clock be the max of the two timestamps (its own and the received), but don't forget that we also need to increment by one. That means if we receive a message with a higher timestamp than our own, we want to change our own to 1  + the message timestamp.
