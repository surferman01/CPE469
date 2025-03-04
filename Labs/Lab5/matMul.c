// #include "mpi.h"
#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>

#define MAX_VAL 1000
#define N 1024
#define CHUNK_SIZE 128

void display_matrix(double **matrix, int size)
{
    for (int i = 0; i < size; i++)
    {
        printf("\t");
        for (int j = 0; j < size; j++)
        {
            printf("%f ", matrix[i][j]);
        }
        printf("\n");
    }
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
            printf("mtx[%d][%d]: %f ", i, j, x);
        }
        printf("\n");
    }
    printf("\n");
    return mtx;
}

double** seqMatMul(double** mtx1, double** mtx2, int size) {
    // allocate space for output mtx
    double** out = (double**)malloc(size * sizeof(double*));
    for (int i = 0; i < size; i++) {
        out[i] = (double*)malloc(size * sizeof(double));
    }

    // perform the matmul (seq)
    for (int i = 0; i < size; i++) {
        for (int j = 0; j < size; j++) {
            for (int k = 0; k < size; k++)  {
                out[i][j] += mtx1[i][k] * mtx2[k][j];
            }
        }
    }
    return out;
}

int main(int argc, char *argv[])
{
    int myid, numprocs, i;
    double **mtx1, **mtx2, **outDist, **outSeq;
    double *row;
    int size = 5;   // just made it a small size here

    mtx1 = genMtx(size);
    mtx2 = genMtx(size);

    outSeq = seqMatMul(mtx1, mtx2, size);

    display_matrix(outSeq, size);

    // standard init for MPI
    MPI_Init(&argc, &argv);
    // numprocs is now the # of spawned procs available
    MPI_Comm_size(MPI_COMM_WORLD, &numprocs);
    // myid is now the rank/id of that specific process
    MPI_Comm_rank(MPI_COMM_WORLD, &myid);

    // scatter rows 
    // so break the mtx1 into the 'row' variable
    // each thread has its own 'row'

    MPI_Scatter(mtx1, size, MPI_DOUBLE, row, size, MPI_DOUBLE, 0, MPI_COMM_WORLD);
    // broadcast entire second mtx
    // i dont think i need to do this since all processes
    // spawn from this one (already have mtx2)
    // MPI_Bcast(mtx2, size*size, MPI_DOUBLE, 0, MPI_COMM_WORLD);

    // many threads spawned here (like fork) where myid is like pid
    if (myid != 0) {
        double* temp = (double*)malloc(size * sizeof(double));
        for (int j = 0; j < size; j++) {
            for (int k = 0; k < size; k++) {
                // use 'myid' since each row that gets scattered
                // allow the threads to operate on the row 
                // that corresponds to its id
                temp[j] += row[k] * mtx2[k][j];
            }
        }
        // gather all the row multiplications into outDist (the output mtx)
        MPI_GATHER(temp, size, MPI_DOUBLE, outDist, size, MPI_DOUBLE, 0, MPI_COMM_WORLD);
    }

    MPI_BARRIER(MPI_COMM_WORLD);
    

//     while (1)
//     {
//         if (myid == 0)
//         {
//             // Generate Matrices
//             srand(time(NULL));
//             X = genMtx(N);
//             srand(time(NULL) ^ 0xDEADBEEF);
//             Y = genMtx(N);

//             // Display input matrices
//             printf("Matrix X:\n");
//             display_matrix(X, N);
//             printf("\nMatrix Y:\n");
//             display_matrix(Y, N);
//         }

//         // Scatter matrices across nodes
//         // MPI_Scatter(matrix, N * chunk_size, MPI_INT, local_matrix, 800 * chunk_size, MPI_INT, 0, MPI_COMM_WORLD); // TODO

//         // Matrix mult

//         // Gather result

//         if (myid == 0)
//         {
//             // Display Result Matrix
//             printf("\nResult Matrix:\n");
//             display_matrix(res, N);
//         }
//     }
//     // MPI_Finalize();
    return 0;
}
