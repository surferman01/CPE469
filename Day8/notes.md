# Day 8 – Feb 4, 2025

## CAP Theorem

- Consistency
  - Every read receives the most recent write or an error
- Availability
  - Every request receives a (non-error) response, w/o guarantee that it contains most recent write
- Partition Tolerance
  - System continues to operate despite an arbitrary number of messages being dropped/delayed by network

## Clock

- Clock Skew
  - distance
- Clock drift
  - difference in oscillation speed

### Clock Sync

Christian's Algorithm

- Let clients contact time server
  - There will be delay, so estimate the delay

NTP

- Network Time Protocol
- Each client sync with parent
  - organized as a tree

Lamport Time Stamps (Logical Clocks)

- Doesn't concern so much with the time
- Care more about order of events

- if a system with a clock ahead of another sends a message to the other, then the one that is behind (receiving the message) will update its clock to be ahead of the sending clock
