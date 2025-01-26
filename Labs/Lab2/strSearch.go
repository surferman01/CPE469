package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Test struct {
	run      int
	name     string
	seqVal   int
	distVal  int
	seqTime  time.Duration
	distTime time.Duration
}

func main() {

	// --------------------
	// EDIT VALUES HERE
	amount := 1
	text := []string{"complete_sherlock.txt", "input"}
	s := "The"
	chunkSize := 20

	// --------------------

	if len(os.Args) > 1 {
		bestSize := 5
		var bestSpeedUp float64 = 1
		for chunkSize := 5; chunkSize < 31; chunkSize++ {
			fmt.Printf("Chunk Size: %d\n", chunkSize)

			out := run_tests(amount, text, s, chunkSize)

			fmt.Println("seq:", out[0].seqTime, "dist:", out[0].distTime, "seqVal:", out[0].seqVal, "distVal: ", out[0].distVal)

			var currTime float64 = float64(time.Duration(out[0].seqTime)) / float64(time.Duration(out[0].distTime))
			if currTime > bestSpeedUp {
				bestSpeedUp = float64(currTime)
				bestSize = chunkSize
			}
		}

		fmt.Printf("\nBest chunk size: %d\n", bestSize)
	} else {
		out := run_tests(amount, text, s, chunkSize)

		fmt.Println("seq:", out[0].seqTime, "dist:", out[0].distTime, "seqVal:", out[0].seqVal, "distVal: ", out[0].distVal)
	}
}

func seq_strSearch(name, s string) int {
	data, err := os.Open(name)
	// fmt.Println(os.ReadFile("input"))
	if err != nil {
		fmt.Println("Error reading input")
		return 1
	}
	defer data.Close()

	counter := 0
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()
		out := strings.Split(line, " ")
		for i := 0; i < len(out); i++ {
			if out[i] == s {
				counter++
			}
		}
	}
	return counter
}

func dist_strSearch(name, s string, chunkSize int) int {
	// Open the file
	data, err := os.Open(name)
	if err != nil {
		fmt.Println("Error reading input")
		return 0
	}
	defer data.Close()

	// Wait group for goroutines processing file chunks
	var wg sync.WaitGroup

	// Wait group for accumulator goroutine
	var accWG sync.WaitGroup

	// Start a goroutine to accumulate results
	total := 0
	buff := make(chan int, 1000)
	accWG.Add(1)
	go func() {
		defer accWG.Done()
		for count := range buff {
			total += count
		}
	}()

	// Read file line by line
	scanner := bufio.NewScanner(data)
	lines := make([]string, 0, chunkSize)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())

		// Collect chunks
		if len(lines) == chunkSize {
			wg.Add(1)
			chunk := append([]string{}, lines...) // Copy lines buffer
			lines = lines[:0]                     // Reset lines buffer

			// Search chunk for string
			go func(chunk []string) {
				defer wg.Done()
				counter := 0
				for _, line := range chunk {
					for _, word := range strings.Fields(line) {
						if word == s {
							counter++
						}
					}
				}
				// Only send data if str is found
				if counter > 0 {
					buff <- counter
				}
			}(chunk)
		}
	}

	// Edge case for any remaining lines in the buffer
	if len(lines) > 0 {
		wg.Add(1)
		chunk := append([]string{}, lines...)
		go func(chunk []string) {
			defer wg.Done()
			counter := 0
			for _, line := range chunk {
				for _, word := range strings.Fields(line) {
					if word == s {
						counter++
					}
				}
			}
			buff <- counter
		}(chunk)
	}

	// Wait for string search goroutines
	wg.Wait()

	// Close the channel
	close(buff)

	// Wait for the accumulating goroutine
	accWG.Wait()

	return total
}

func run_tests(amount int, text []string, s string, chunkSize int) []Test {
	out := make([]Test, amount)

	if amount > len(text) {
		fmt.Println("The amount needs to be less than the # of texts")
		return nil
	}

	for run := 0; run < amount; run++ {
		// run is just the run number
		out[run].run = run

		fmt.Println("run:", run)

		// name is the text file name being processed
		name := text[run]
		out[run].name = name

		beforeSeq := time.Now()
		seq_srch := seq_strSearch(name, s)
		sinceSeq := time.Since(beforeSeq)

		out[run].seqTime = sinceSeq

		beforeDist := time.Now()
		dist_srch := dist_strSearch(name, s, chunkSize)
		sinceDist := time.Since(beforeDist)

		out[run].distTime = sinceDist

		check := check_result(seq_srch, dist_srch)
		if check != true {
			fmt.Println("The search trials are not equal for run:", run, name)
			// for testing, the below is commented
			// return nil
		}
		// val is the number of findings
		out[run].distVal = dist_srch
		out[run].seqVal = seq_srch
	}
	return out
}

func check_result(out1, out2 int) bool {
	if out1 == out2 {
		return true
	}
	return false
}
