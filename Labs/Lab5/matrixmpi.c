#include <mpi.h>
#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <sys/time.h>

#define MAX_VAL 10
#define N 1024
#define THRESHOLD 10

double **generate_matrix(int n);
void display_matrix(double *M);
void multiply_matrices(double *A, double *B, double *C, int rows_per_proc, int n);
int comp_matrices(double *A, double *B);

int main(int argc, char *argv[])
{
    int rank, size;
    struct timeval start, seq_end, dist_end;
    MPI_Init(&argc, &argv);
    MPI_Comm_rank(MPI_COMM_WORLD, &rank);
    MPI_Comm_size(MPI_COMM_WORLD, &size);

    int rows_per_proc = N / size; // Assume N is divisible by size
    double A[N * N], B[N * N], C[N * N], D[N * N];
    double local_A[rows_per_proc * N], local_C[rows_per_proc * N];

    if (rank == 0)
    {
        // Initialize matrices A and B
        for (int i = 0; i < N * N; i++)
        {
            A[i] = i + 1;
            B[i] = (i + 1) % N;
        }

        if (N <= THRESHOLD)
        {
            // Display Inputs
            printf("Matrix A: \n");
            display_matrix(A);
            printf("Matrix B: \n");
            display_matrix(B);
        }

        // Sequential Matrix Mult
        gettimeofday(&start, NULL);
        multiply_matrices(A, B, D, N, N);
        gettimeofday(&seq_end, NULL);
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
        gettimeofday(&dist_end, NULL);

        if (N <= THRESHOLD)
        {
            // Print result matrix C
            printf("Result:\n");
            display_matrix(C);
        }

        // Display statistics
        if (comp_matrices(C, D) < 0)
        {
            printf("Sequential and distributed solutions DON'T match\n");
        }
        else
        {
            printf("Sequential and distributed solutions DO match\n");
        }

        double seq_time = (seq_end.tv_sec - start.tv_sec) + (seq_end.tv_usec - start.tv_usec) / 1000000.0;
        double dist_time = (dist_end.tv_sec - seq_end.tv_sec) + (dist_end.tv_usec - seq_end.tv_usec) / 1000000.0;

        printf("Sequential Time: %f\n", seq_time);
        printf("Distributed Time: %f\n", dist_time);
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

void display_matrix(double *M)
{
    for (int i = 0; i < N; i++)
    {
        printf("\t");
        for (int j = 0; j < N; j++)
        {
            printf("%f ", M[i * N + j]);
        }
        printf("\n");
    }
}

int comp_matrices(double *A, double *B)
{
    for (int i = 0; i < N * N; i++)
    {
        if (A[i] != B[i])
        {
            return -1;
        }
    }

    return 1;
}
