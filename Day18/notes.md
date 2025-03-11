# Day 18 (3/11/2025)

## Convolution / Gaub-Seidel

### Convolution

- when performed on an array, generate a new one based on inpit

### Gaub-Seidel

- when performed on an array, modify the current array (less space complexity)
- in the convolution, use a stencil that considers the cardinal direction values also in the array
- then distribute the conv by the diagonals

## MPI Sending

MPI_Bsend

- buffer (blocking)
- sending operation locally blocks
- buffer is user provided

MPI_Ssend

- synchronous
- sending operation will return only after destination process has started receiving message
- global blocking operation

MPI_Rsend

- ready
- send will only succeed if a matching receive operation has been started on the receiver

MPI_Isend

- fully async
- can put I in front of any send to make it go immediately

## N-Body Problems

- use one body per node (conv style)
- or use tiling
