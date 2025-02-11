# Day 10 (2/11/2025)

## Post Midterm Housekeeping

### Solution to General Problem

- No single solution that always works
- Item Potency
  - If you want to add to a database, ensure it is added only ONCE

## Recap – 3 components for DS

Membership protocol
Consensus
Failure Detection

## DynamoDB (noSQL, no rdbm, no tables)

### RDBM: Relational database manager

- all information stored on tables
- attributes = column headers
- rows = tuples

Don't store all information in a single relation (row)

- creates redundant data
- need null values if elements missing

### SQL

- say "what to do" instead of "how to do it"
- database management system handles query
  - optimized by RDBM – query optimization

Cartesian Product

- for each row of one table, give it an entire other table for its data

Natural Join

- match tuples w same values for all common attrs and retain only one copy of each common column

Transactions

- Atomic
  - Either fully executed (comitted) or rolled back

### Consistent Hashing

Hashing recap

- Map ID to some value
- simplest implementation
  - id -> value
  - id % # values = key
  - use % to map keys to values

Consistent Hashing

- use a circle of nodes where each node holds a range of keys
- let them overlap such that each key has 3 replicas
