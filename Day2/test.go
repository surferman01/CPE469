package main

import (
	"fmt"
	"math"
	"os"
)

type employee struct {
	ID   int
	Name string
}

func main() {
	x := 1
	p := &x
	fmt.Println(*p) // --> 1
	fmt.Println(p)  // --> 0xaddr

	s := "test"
	var y int = 10
	fmt.Println(y) // --> 10
	fmt.Println(s) // --> test

	const pi = 3.14

	var s2, sep string
	for i := 1; i < len(os.Args); i++ {
		s2 += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s2)

	// concurrent programming

	t := make(chan bool)
	// go timeout(t)

	ch := make(chan string)
	// go readword(ch)

	// select {...}

	ch <- x		// send statement
	x =<-ch 	// receive
	<- ch		// receive, result discarded
	close()		// no more communication
	
	ch2 := make(chan int, 3) 	// buffered channel with capacity 3

}

// func name (arg-list) (return-list types) {}
func hyp(x, y float64) float64 {
	return math.Sqrt(x*y + y*y)
}

func square() func() int {
	var x int
	return func() int {
		x++
		return x * x
	}
}

func sum(vals ...int) int {
	total := 0
	for _, vals := range vals {
		total += vals
	}
	return total
}
