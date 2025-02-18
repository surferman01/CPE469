# Day 12, (2/18/2025)

## Transactions

- transaction = unit of program execution that accesses / updates data
- must see consistent database
- during transaction, database temporarily inconsistent
- when transaction completes (is comitted), database must be consistent
- after transaction, data must persist even if system fail
- must support parallel transactions
- hardest problems
  - concurrency of transactions
  - failures

### ACID Properties

- Atomicity: either all transactions are properly in database or none
- Consistency: execution of transaction in isolation must preserve consistency of database
- Isolation: each transaction must be unaware of other transactions. Intermediate transactions must be hidden from other concurrent transactions
- Durability: after transaction, changes must persist in db even if sys fail

### Log Based Recovery

- use checkpoints so you know how far back to go to get everything updated

### Schedules

- schedule: sequence of instructions that specify chronological order in which instructions of concurrent transactions are executed
  - schedule for set of transactions must have all instructions of those transactions
  - must preserve order in which instructions appear in each individual transaction
- transaction the successfully completes will have a commit at the end
- transaction that fails will have an abort at the end

### Serializability

- serial: instructions one after another, one at a time
- precedence graph: direct graph where verticies are transaction names
  - draw an arc if two transactions conflict

## Document Collaboration

### Operational Transform

- server is responsible for merging changes and sending back to clients

### Conflict free replicated data type (CRDT)

- data structure that is replicated across multiple computers in a network

features:

- application can update replica independently, concurrently and w/o coordinating w other replicas
- an algorithm (embedded in the data type itself) automatically resolves inconsistencies
- although replicas may temporarily have diff states, they are guaranteed to eventually converge
