#include "mpi.h"
#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>

#define MAX_VAL 1000
#define N 1024
#define CHUNK_SIZE 128

double **generate_matrix(int n);
void display_matrix(double **matrix, int n);

int main(int argc, char *argv[])
{
    int myid, numprocs, i;
    double **res, **X, **Y;

    MPI_Init(&argc, &argv);
    MPI_Comm_size(MPI_COMM_WORLD, &numprocs);
    MPI_Comm_rank(MPI_COMM_WORLD, &myid);

    while (1)
    {
        if (myid == 0)
        {
            // Generate Matrices
            srand(time(NULL));
            X = generate_matrix(N);
            srand(time(NULL) ^ 0xDEADBEEF);
            Y = generate_matrix(N);

            // Display input matrices
            printf("Matrix X:\n");
            display_matrix(X, N);
            printf("\nMatrix Y:\n");
            display_matrix(Y, N);
        }

        // Scatter matrices across nodes
        MPI_Scatter(matrix, N * chunk_size, MPI_INT, local_matrix, 800 * chunk_size, MPI_INT, 0, MPI_COMM_WORLD); // TODO

        // Matrix mult

        // Gather result

        if (myid == 0)
        {
            // Display Result Matrix
            printf("\nResult Matrix:\n");
            display_matrix(res, N);
        }
    }
    MPI_Finalize();
    return 0;
}


double** genMtx(int size) {
    double** mtx = (double**)malloc(size * sizeof(double*));
    for (int i = 0; i < size; i++) {
        mtx[i] = (double*)malloc(size * sizeof(double));
    }
    for (int i = 0; i < size; i++) {
        for (int j = 0; j < size; j++) {
            double a = 10.0;
            double x = ((double)rand()/(double)(RAND_MAX)) * a;
            mtx[i][j] = x;
            printf("mtx[%d][%d]: %f\n", i, j, x);
        }
    }
    return mtx;
}