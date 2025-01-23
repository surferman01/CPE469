package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Test struct {
	run      int
	name     string
	val      int
	seqTime  time.Duration
	distTime time.Duration
}

func main() {

	// --------------------
	// EDIT VALUES HERE
	amount := 1
	text := []string{"complete_sherlock.txt", "input"}
	s := "The"
	// --------------------

	out := run_tests(amount, text, s)

	fmt.Println("seq:", out[0].seqTime, "dist:", out[0].distTime, "val:", out[0].val)
	// s := "test"

	// data, err := os.ReadFile("input")
	// if err != nil {
	// 	fmt.Println("Error reading input")
	// 	return
	// }

	// asStr := string(data)

	// out := strings.Split(asStr, " ")

	// counter := 0
	// for i := 0; i < len(out); i++ {
	// 	if s == out[i] {
	// 		counter++
	// 	}
	// 	fmt.Println(s, "S HERE")
	// 	fmt.Println(out[i])
	// }

	// fmt.Println(counter)

	// fmt.Println(out)

	// BUFIO / SCANNER implementation

	// fmt.Println(seq_strSearch("complete_sherlock.txt", "test"))
	// data, err := os.Open("input")
	// // fmt.Println(os.ReadFile("input"))
	// if err != nil {
	// 	fmt.Println("Error reading input")
	// 	return
	// }
	// defer data.Close()

	// counter := 0
	// s := "test"
	// scanner := bufio.NewScanner(data)
	// for scanner.Scan() {
	// 	line := scanner.Text()
	// 	out := strings.Split(line, " ")
	// 	for i := 0; i < len(out); i++ {
	// 		fmt.Println(out[i])
	// 		if out[i] == s {
	// 			counter++
	// 		}
	// 	}
	// }
	// fmt.Println(counter)
	// fmt.Println(s)
	// for i := 0; i < len(data.); i++ {

	// }
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

func dist_strSearch(name, s string) int {

	// TODO
	return 0
}

func run_tests(amount int, text []string, s string) []Test {
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
		dist_srch := dist_strSearch(name, s)
		sinceDist := time.Since(beforeDist)

		out[run].distTime = sinceDist

		check := check_result(seq_srch, dist_srch)
		if check != true {
			fmt.Println("The search trials are not equal for run:", run, name)
			// for testing, the below is commented
			// return nil
		}
		// val is the number of findings
		out[run].val = dist_srch
	}
	return out
}

func check_result(out1, out2 int) bool {
	if out1 == out2 {
		return true
	}
	return false
}
