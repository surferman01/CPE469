# Day 4, (1/16/2025)

## Error Handling

- We can handle in 2 main ways
  - Let the goroutine fill up a data structure with all info (data + error) then afterwards process it
  - Create a WaitGroup
    - like a block that syncs the program.
    - see test.go
    - need to use 'defer' keyword
      - this makes the function call after the 'defer' keyword execute after the main program exits (like a 'finally' in try/catch)

## Race Conditions

- Go has built in race detector
  - just add -race flag to go run/build

## MapReduce

- Map: (input shard) -> intermediate(key/value pairs)
  - read much data and extract something relevant from each record
  - Steps
    1) partition input data into M shards
    2) discard unnecessary data and generate (key, value) sets
    3) framewor groups together all intermediate values with same intermediate key and pass to Reduce function
- Reduce: intermediate(key/value pairs) -> result files
  - summarize/filter/transform data
  - merge values together to form smaller set of values
  - use partitioning function to split maps into R partitions

### Steps

1) See slides
