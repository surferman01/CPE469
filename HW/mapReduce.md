# MapReduce Notes

- User specified `map` function that processes key/value pairs to generate intermediate key/value pairs
- reduce function that merges the intermediate values that share a common key
- this tech enables distributed computing w/o requiring underlying knowledge about it

## Master Worker

- one master for a fleet of mapreduce workers
- periodically pings the workers, if one goes offline its marked as failure
  - its work is scrapped and the data is sent to another worker to be processed
    - this is mainly because the intermediate data is stored locally on that worker, and we no longer can access it
- if the sole master fails, it will cease mapreduces operations
- mappers generate temp files that get sent to master once it completes
- those temp files are sent to reducers that operate on them until one final cut is created
- that is renamed to the expected output file name and sent to the master

- subdivide map and reduce into M and R tasks
  - these should be much more tasks than the number of workers for parallelism and failure integrity
- the master needs to make O(M + R) scheduling decisions
- it needs O(M * R) space though
  - luckily, one task can be defined by only a byte per pair!
- size ex:
  - M = 200 000
  - R = 5 000
  - workers = 2 000

- as a MapReduce call is soon to end, spawn extra 'backup' tasks so if one processes slow another can do it faster and send result to master
- master accepts the primary one OR the backup one (whichever is first)
- in their example, backup produces a 44% boost in speed

- the mapReduce output is always ordered

- if an error is detected on some map or reduce more than once, it is skipped and flagged

- there are counters nested within mapreduce that can be used for sanity checks
