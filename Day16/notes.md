# Day 16 (3/4/2025)

## Median of an array

- emit some value to all nodes
- have each node return the values below/above
- remove one of the groups
- repeat until you find the median

## How to sort a massive array

- split data into chunks and sort on each computer
- now need to m/erge sorted arrays
  - not enough space on one computer for holding multiple chunks
  - do odd_even sort

### Odd_Even sort

- compare adjacent nodes and merge
- they are already sorted, so just do a regular merge
- if one value is greater than the other, then move accordingly
