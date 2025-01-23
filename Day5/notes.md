# Day 5 – 1/23/2025

## Communications for DS

- Needs to be reliable + scalable
- Traditionally based on low level message passing
  - Remote Procedure Calls (RPC)
  - Messsage Oriented
  - Multicasting

### Remote Procedure Calls (RPC)

Make function (procedure) calls from a master to a remote. Then, the remote sends back the result to the master.

- Utilizes socket for communication
  - thru TCP

### Multicast Communication

Sending data to multiple receivers

- Tree based
  - One node talks to two
  - Parent to children communication
    - If one node goes offline, then whole branch offline
  - Use ack – acknowledge to determine online/not

### Gossip

Periodically transmit data to random target. Other nodes do same after receiving multicast.

The purpose is to get one message to all the nodes. One node sends message to x others, when they recieve it, they send it to x others until everyone has it.

- Built on TCP
- Random / Redundant
- Sends more messages than other communications
  - Push: send message when you have message
  - Pull: see if nodes around you have received a message recently

## Failure Detection

### Types

Byzantine / Arbitrary failures:

A server sending contrary messages to different servers

- Authentication
  - Server realizes its sending false information, but cannot alter ripple of information from other servers
- Performance
  - Server takes too long to read / send messages
  - Omission
    - server replies infinitely late
  - Crash
    - after omission, stop responding
  - Crash-stop
    - process halts and stops responding

### Failure Detector

- Completeness
  - Node fails, then detect it
  - must be detected by non-faulty process
  - must be accurate

### System Models

Synchronous Distributed System

- messages recieved within bounded time
- ex)
  - multiprocessor systems

Asynchronous DS

- No bounds on transmission delay
- ex)
  - Internet, wireless networks, datacenters

### Membership

- Each process has a partial table of other processes in group
- Updated with changes from new members or leaving members
- In gossip, one node has table of all other nodes

### Heartbeet Protocol

One master has a sequence number that sends heartbeats to node and wait for response. Using the sequence number, it can determine the delay for the responses.

Centralized Heartbeat

- one master sends / receives heartbeats from all

Ring Heartbeat

- sequentially go thru nodes like linked list

All-to-All Heartbeat

- as name implies

Gossip Style

- periodically gossip to random node your heartbeat
- pass it around and update tables
