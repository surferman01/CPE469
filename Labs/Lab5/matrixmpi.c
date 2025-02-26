#include "mpi.h"
#include <stdio.h>
#include <math.h>

int main(int argc, char *argv[])
{
    int n, myid, numprocs, i;
    double sum, x;

    MPI_Init(&argc, &argv);
    MPI_Comm_size(MPI_COMM_WORLD, &numprocs);
    MPI_Comm_rank(MPI_COMM_WORLD, &myid);

    while (1)
    {
        if (myid == 0)
        {
            // stdio operations
            printf("Enter num intervals (0 quits): ");
            scanf("%d", &n);
        }

        MPI_Bcast(&n, 1, MPI_INT, 0, MPI_COMM_WORLD);

        if (n == 0)
        {
            break;
        }
        else
        {
            // Matrix mult here

            if (myid == 0)
            {
                // Display Matrix
            }
        }
    }
    MPI_Finalize();
    return 0;
}