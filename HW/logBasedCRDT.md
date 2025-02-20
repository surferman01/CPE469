# Log-Based CRDT for Edge Applications Notes (2/19/2025) (https://ieeexplore.ieee.org/document/9946360)

## CRDT

- Conflict-free Replicated Data Types

## Problem

- many CRDT formulations do not support operation reversal
  - reverting a datatype to a previous state for operations that are not executed
- cannot tolerate out of order operations

## Solution

- use distributed logs
- since append only, it reduces coordination required to merge inconsistent replicas
- LSCRDT

## Metrics

- latency
- throughput

- measured using 3 CRDT data types
  1. register
  2. counter
  3. set

## Types of CRDTs

1. state based
    - operation located at local replica
    - periodically emit data to other replicas for consistency
    - maybe lots of overhead passing entire state to other replicas
2. operation based
    - each operation happend locally, but also is asynchronously propagated to others

## LSCRDT

- capable of tolerating unreliable networking + manage arbitrary datatypes
- each operation logs executed operation WITH unique verison stamp (for duplicate detection) in causal order
- maintains causal order by ensuring log order is maintained during reads
- logs are immutable
- operation based
  - so that each implementation is fair
  - usually CRDTs have some merge protocol
    - using consistent order of operations, no merge protocol necessarily required
    - thus, any datatypes can be used

## Implementation

- if replicas diverge, perform merge steps
- merge steps are always between 2 replicas

- One replica becomes READER
  - reader reads logs from source
- One replica becomes SOURCE

- goal:
  - for reader to identify operations "unknown" to it
  - reader will rollback its operations (undo) until it has all the latest operations
  - then it will re-execute all the operations until up to date

- use timestamps to determine 'correct' ordering

## Log Requirements

1. create logs with given name
2. write to a specified log and get the sequence number corresponding to the write on success
3. read from a specified log at a given sequence number
4. retrieve the latest sequence number of a specified log
5. trim the log up to a specified sequence number

- use Lamport timestamps

## Register

- register = data type
- maintains SINGLE value
- 2 ops
    1. assign (set)
    2. retreive (get)

- use OpLog to store all update ops (assigns)
- have ONE OpLog per replica

- Replicas can read each others OpLog

- OpLog tuple: (vs, op, val)
- vs = version stamp
- op = type of update (for register, theres only one: assign)
- val = operand of op

## Counter

counter = datatype

- 2 ops
    1. increment (inc)
    2. decrement (dec)

- use an OpLog to see prev value, new value, change in val, and time

- OpLog tuple: (vs, op, amt, val)
- vs = version stamp
- op = type of update (inc or dec)
- amt = how much to inc or dec by
- val = new val

## Set

set = datatype

- 4 ops
  - update operands
        1. add
        2. remove
  - read operands
        1. in
        2. all

- each set = collection of elements
- to reconstruct a set, we need to examine OpLog

## Merge Step

- basically like RAFT log replication
- go back in history until you see all the 'unknown operations'
- then start appending to your OpLog (and performing) the operations in the lamport order

reader performs 2 tasks

1. conflict detection:
    - see if the source has 'unknown' operations compared to reader
    - if the source has missing data from the reader, do NOTHING
2. conflict resolution
    - reader rollsback their operations, and executes operations with newly learned operations mixed into rolled-back operations

### KnowledgeLogs

- KnowledgeLog = map of last observed version stamp from each replica to a sequence number in its OpLog
- use a KnowledgeLog for each replica
- tuple: (vs, op_seq)
  - vs = version stamp
  - op_seq = sequence number of OpLog where the operation with vs was first appended
- its like the last synced update with each replica (starting point when performing merge step)
- use this so that we dont need to reference an entire OpLog when we read from a source
- also cache the important data so its easily looked up

### Conflict Detection

- scan until you find the earliest operation in the OpLog that is unknown to the reader
- compare vs values and ensure log matching between them

### Conflict Resolution

- roll back OpLog of the reader to the earliest unknown operation
- dont go past the earliest version stamp that has been synced already
- replay operations at reader to match merge order

### Bounded Log Sizes

- very confusing section...

## Evaluation

- how much more time since using logs
- effect of logs on scalability
- how time-consuming is versioned reads compared to reading latest value

- times:
  - slower set read/writes
  - this is because the delta-CRDT uses a 2P set where it can just reference logs marked as ready to be removed for quick comparisons while LSCRDT needs to go back through the OpLog until it has seen the earliest unknown operation. Then it needs to rebuild from there

- for edge applications, LSCRDT is a better options since these applications typically are write-heavy, more updates

### Scalability

- `Note that LSCRDT set has a higher read latency but a lower write latency than δ-CRDT set. The lower write latency along with the efficient lock-free merge step of LSCRDT results in a higher throughput in LSCRDT set for workloads with a high volume of updates.` (article)

### Versioned Reads

- maintaining version history is unique to LSCRDTs, not available in CRDTs
- the versioned reading takes longer than that of CRDTs, but the write is faster (primary purpose)

## Conclusion

- LSCRDT is very useful for supporing any arbitrary datatype of storage in the logs
- its the first CRDT that tracks version histories of data structures
- very useful for update-heavy workloads
- improve throughput by 1.8x
