# Day 14 (2/25/2025)

## MPI

goal:

provide widely used standard for writing message passing programs

### Compiling

mpicc -g -Wall -o mpi_test mpi_test.c

### Executing

mpiexec -n [# nodes] ./mpi_test

### Reduce

MPI_Allreduce

- if you distribute processing/data, you need to reconcile the data into one sum @ end
- to do this, use butterfly reduction
  - pair off and send your value to partner
  - now you both have same sum
  - then, pair with another pair and send sums
  - now all 4 have same sum
  - continue until consensus

### Broadcast

- data in one process is sent to all processes
- you pass your data + you name
- that means it is the receiving nodes responsibility to act on the data if it needs
  - if the name of the broadcaster is the one you care about, then act

### Scattering

- for distributing data (in an array) to many nodes

### Gathering

- collecting all the data at the end (after reduce)
