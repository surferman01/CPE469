package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/wcharczuk/go-chart/v2"
)

// matrix:
// [1, 2, 3]
// [4, 5, 6]
// [7, 8, 9]

func main() {
	// --------------------
	// EDIT VALUES HERE
	amount := 1
	increment := 0
	start := 1000
	// --------------------

	printResults(runTests(amount, increment, start))
	// runTests(1, 0, 1500)
}

func genMtx(size int) [][]float64 {
	mtx := make([][]float64, size)
	for i := range size {
		mtx[i] = make([]float64, size)
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			mtx[i][j] = rand.Float64()
			// fmt.Println(mtx[i][j], "i: ", i, "j: ", j)
		}
	}
	return mtx
}

func matMulSeq(m1, m2 [][]float64) [][]float64 {
	out := make([][]float64, len(m1))
	// fmt.Println(len(m1))

	for i := range out {
		out[i] = make([]float64, len(m1[0]))
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

func matMulDist(m1, m2 [][]float64) [][]float64 {
	out := make([][]float64, len(m1))
	for i := range out {
		out[i] = make([]float64, len(m1))
	}
	// important to use an index like this since
	// we dont know when the math will finish
	// so knowing/passing the index guarantees
	// it will be in the right place later
	element := make(chan struct {
		i    int
		data []float64
	}, len(m1))

	for i := 0; i < len(m1); i++ {
		// element := make(chan []float64)
		// go getMatMulRow(m1, m2, i, element)

		go func(i int) {
			out := make([]float64, len(m1))
			for j := 0; j < len(m1); j++ {
				for k := 0; k < len(m1); k++ {
					out[j] += m1[i][k] * m2[k][j]
				}
			}
			element <- struct {
				i    int
				data []float64
			}{i: i, data: out}
		}(i)
	}
	for i := 0; i < len(m1); i++ {
		temp := <-element
		out[temp.i] = temp.data
	}
	return out
}

// func getMatMulRow(m1, m2 [][]float64, i int, element chan<- []float64) {
// 	out := make([]float64, len(m1))
// 	for j := 0; j < len(m1); j++ {
// 		for k := 0; k < len(m1); k++ {
// 			out[j] += m1[i][k] * m2[k][j]
// 		}
// 	}
// 	element <- out
// }

func checkMatMul(out1, out2 [][]float64) bool {
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

func runTests(amount, increment, start int) []struct {
	run      int
	size     int
	seqTime  time.Duration
	distTime time.Duration
} {
	out := make([]struct {
		run      int
		size     int
		seqTime  time.Duration
		distTime time.Duration
	}, amount+1)

	for run := 1; run < amount+1; run++ {
		size := start + increment*run
		mtx1 := genMtx(size)
		mtx2 := genMtx(size)

		fmt.Println("run:", run, "size:", len(mtx1))
		// fmt.Println(len(mtx2))
		// fmt.Println("first line:", matrix2[0][0])
		beforeSeq := time.Now()
		mtxoutSeq := matMulSeq(mtx1, mtx2)
		afterSeq := time.Since(beforeSeq)

		fmt.Println("sequential done:", afterSeq)

		beforeDist := time.Now()
		mtxoutDist := matMulDist(mtx1, mtx2)
		afterDist := time.Since(beforeDist)

		fmt.Println("distributed done:", afterDist)
		fmt.Println("check?: ", checkMatMul(mtxoutDist, mtxoutSeq))
		fmt.Println()

		out[run] = struct {
			run      int
			size     int
			seqTime  time.Duration
			distTime time.Duration
		}{run, size, afterSeq, afterDist}
	}
	fmt.Println("Done.")
	return out
}

func printResults(result []struct {
	run      int
	size     int
	seqTime  time.Duration
	distTime time.Duration
}) {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "Matrix Size",
		},
		YAxis: chart.YAxis{
			Name: "Time (ns)",
		},
		Series: []chart.Series{},
	}

	seqSeries := chart.ContinuousSeries{
		Name:    "Sequential Time",
		XValues: []float64{},
		YValues: []float64{},
	}

	distSeries := chart.ContinuousSeries{
		Name:    "Distributed Time",
		XValues: []float64{},
		YValues: []float64{},
	}

	for _, res := range result {
		seqSeries.XValues = append(seqSeries.XValues, float64(res.size))
		seqSeries.YValues = append(seqSeries.YValues, float64(res.seqTime.Nanoseconds()))
		distSeries.XValues = append(distSeries.XValues, float64(res.size))
		distSeries.YValues = append(distSeries.YValues, float64(res.distTime.Nanoseconds()))
	}

	graph.Series = append(graph.Series, seqSeries, distSeries)

	f, _ := os.Create("output.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
}
