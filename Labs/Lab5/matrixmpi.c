#include <mpi.h>
#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>

#define MAX_VAL 10
#define N 1024

double **generate_matrix(int n);
void display_matrix(double **matrix, int n);
void multiply_matrices(double *A, double *B, double *C, int rows_per_proc, int n);

int main(int argc, char *argv[])
{
    int rank, size;
    MPI_Init(&argc, &argv);
    MPI_Comm_rank(MPI_COMM_WORLD, &rank);
    MPI_Comm_size(MPI_COMM_WORLD, &size);

    int rows_per_proc = N / size; // Assume N is divisible by size
    double A[N * N], B[N * N], C[N * N];
    double local_A[rows_per_proc * N], local_C[rows_per_proc * N];

    if (rank == 0)
    {
        // Initialize matrices A and B
        for (int i = 0; i < N * N; i++)
        {
            A[i] = i + 1;
            B[i] = (i + 1) % N;
        }
    }

    // Broadcast matrix B to all processes
    MPI_Bcast(B, N * N, MPI_DOUBLE, 0, MPI_COMM_WORLD);

    // Scatter rows of A to all processes
    MPI_Scatter(A, rows_per_proc * N, MPI_DOUBLE, local_A, rows_per_proc * N, MPI_DOUBLE, 0, MPI_COMM_WORLD);

    // Perform local matrix multiplication
    multiply_matrices(local_A, B, local_C, rows_per_proc, N);

    // Gather results from all processes
    MPI_Gather(local_C, rows_per_proc * N, MPI_DOUBLE, C, rows_per_proc * N, MPI_DOUBLE, 0, MPI_COMM_WORLD);

    if (rank == 0)
    {
        // Print result matrix C
        for (int i = 0; i < N; i++)
        {
            for (int j = 0; j < N; j++)
            {
                printf("%f", C[i * N + j]);
            }
            printf("\n");
        }
    }

    MPI_Finalize();
    return 0;
}

// Function to multiply matrices
void multiply_matrices(double *A, double *B, double *C, int rows_per_proc, int n)
{
    for (int i = 0; i < rows_per_proc; i++)
    {
        for (int j = 0; j < n; j++)
        {
            C[i * n + j] = 0;
            for (int k = 0; k < n; k++)
            {
                C[i * n + j] += A[i * n + k] * B[k * n + j];
            }
        }
    }
}
