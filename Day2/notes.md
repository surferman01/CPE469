# Day 2, Jan. 9, 2025

## Go Intro

Functions can have multiple returns

func name (arg-list) (return-list types) {}

### Anonymous Functions

- Functions without names

### Concurrent Programming

- 2 styles
  - Goroutines and channels
  - Shared memory multithreading
- putting 'go' before a function makes it run asynchronously

- Channels
  - make(chan [type])
  - connections between go routines allowing them to send values to each other
  - Unbuffered vs Buffered
    - you can determine the buffer size for a channel 
    - Buffered
      - make(chan int, 3)
      - ^ a channel with buffer size 3
    - Unbuffered
      - make(chan string)

