# Day 3, (1/14/2025)

## Some recap

### Processes

- get own memory space independently

### Threads

- 'lightweight processes'
- use same memory space as parent process
- gets own new stack
- requires 2MB of stack memory
- scheduled by OS -> OS controls context
- since threads create local variables, there is no raised condition possibility

### Go Routine

- starts with 2kB of memory and can grow to 1GB
- have own scheduler -> controlled by software
- no PID for goroutines
- gomaxprocs -> max amount of threads that execute go code simultaneously
- executes on same address space as other goroutines
- prefixing a call with 'go' will make it run in a goroutine

### Channels

- connection between go routines allowing communication
- unbuffered (2 operations, send/receive):
  - ch := make(chan int)
  - ch<-x       // send statement
  - x =<-ch     // receive
  - <-ch        // receive, result discarded
  - close()
- buffered:
  - ch=make(chan int, 3)  // buffer capacity = 3
- receivers block until data exists in chan

### Functions / Notes

- range is like python, but you can supply 2 input vars and get the index + value
  - list := [6]int{1, 2, 4, 8, 16, 32}
  - for i, val := range list {}
