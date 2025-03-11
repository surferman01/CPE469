CC = mpicc
CFLAGS = -Wall

# Targets
all: matMul

matMul: matrixmpi.c
        $(CC) $(CFLAGS) -o matMul matrixmpi.c

sweep:
        rm -f core* *.out

clean:
        rm -f matMul