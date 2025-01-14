package main

import (
	"fmt"
	"math/rand"
	"time"
)

// matrix:
// [1, 2, 3]
// [4, 5, 6]
// [7, 8, 9]

func main() {
	size := 1500
	// rows := make([]float32, size)
	// cols := make([]float32, size)
	// for i := 0; i < size; i++ {
	// 	for j := 0; j < size; j++ {

	// 	}
	// 	rows[i] = rand.Float32()
	// 	cols[i] = rand.Float32()
	// 	fmt.Println(rows[i], cols[i])
	// }

	// matrix1 := make([][]float32, size)
	// matrix2 := make([][]float32, size)
	// for i := range size {
	// 	matrix1[i] = make([]float32, size)
	// 	matrix2[i] = make([]float32, size)
	// }
	// for i := 0; i < size; i++ {
	// 	for j := 0; j < size; j++ {
	// 		matrix1[i][j] = rand.Float32()
	// 		matrix2[i][j] = rand.Float32()
	// 		fmt.Println(matrix1[i][j], "i: ", i, "j: ", j)
	// 	}
	// }

	mtx1 := genMtx(size)
	mtx2 := genMtx(size)
	fmt.Println(len(mtx1))
	// fmt.Println(len(mtx2))
	// fmt.Println("first line:", matrix2[0][0])
	beforeSeq := time.Now()
	mtxoutSeq := matMulSeq(mtx1, mtx2)
	afterSeq := time.Since(beforeSeq)

	fmt.Println("sequential done")

	beforeDist := time.Now()
	mtxoutDist := matMulDist(mtx1, mtx2)
	afterDist := time.Since(beforeDist)

	fmt.Println("distributed done")

	// fmt.Println(mtxoutDist[0][0])

	fmt.Println("sequential time:", afterSeq)
	fmt.Println("dist time:", afterDist)
	fmt.Println("check?: ", checkMatMul(mtxoutDist, mtxoutSeq))

}

func genMtx(size int) [][]float32 {
	mtx := make([][]float32, size)
	for i := range size {
		mtx[i] = make([]float32, size)
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			mtx[i][j] = rand.Float32()
			// fmt.Println(mtx[i][j], "i: ", i, "j: ", j)
		}
	}
	return mtx
}

func matMulSeq(m1, m2 [][]float32) [][]float32 {
	out := make([][]float32, len(m1))
	// fmt.Println(len(m1))

	for i := range out {
		out[i] = make([]float32, len(m1[0]))
	}

	for i := 0; i < len(m1); i++ {
		for j := 0; j < len(m1); j++ {
			for k := 0; k < len(m1[0]); k++ {
				out[i][j] += m1[i][k] * m2[k][j]
			}
		}
	}

	return out
}

func matMulDist(m1, m2 [][]float32) [][]float32 {
	// size := len(m1)
	out := make([][]float32, len(m1))
	// output := make(chan [][]float32, len(m1))
	// fmt.Println(len(m1))

	for i := range out {
		out[i] = make([]float32, len(m1))
	}
	element := make(chan []float32)

	for i := 0; i < len(m1); i++ {
		// element = make(chan float32)
		// go func()
		// element := make(chan []float32)
		go getMatMulRow(m1, m2, i, element)

		out[i] = <-element
		// close(element)

		// for j := 0; j < len(m1); j++ {
		// 	// go func()
		// 	element := make(chan float32)
		// 	// go getMatMulElement(m1, m2, i, j, element)
		// 	// for k := 0; k < len(m1[0]); k++ {
		// 	// 	out[i][j] += m1[i][k] * m2[k][j]
		// 	// }
		// 	out[i][j] =<-element
		// 	close(element)
		// }
	}
	close(element)

	return out
}

// func matMulRow(r1, r2 [][]float32) [][]float32 {

// }

func getMatMulRow(m1, m2 [][]float32, i int, element chan<- []float32) {
	out := make([]float32, len(m1))
	for j := 0; j < len(m1); j++ {
		for k := 0; k < len(m1); k++ {
			out[j] += m1[i][k] * m2[k][j]
		}
	}
	element <- out
}

func checkMatMul(out1, out2 [][]float32) bool {
	size := len(out1)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if out1[i][j] != out2[i][j] {
				fmt.Println(out1[i][j], out2[i][j], i, j)
				return false
			}
		}
	}
	return true
}
