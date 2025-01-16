package main

import (
	"os"
	"sync"
	"thumbnail"
)

var wg sync.WaitGroup // number of working subroutines
func makeThumbnails(filenames <-chan string) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup // number of working subroutines
	for f := range filenames {
		thumb, err := thumbnail.ImageFile(f)
		go func(f string) {
			defer wg.Done() //equivalent to add(-1); defer => dec counter even on error
			filenames, err = thumbnail.ImageFile(f)
			if err != nil {
				info, _ := os.Stat(thumb)
				sizes <- info.Size()
			}
			info, _ := os.Stat(thumb)
			sizes <- info.Size()
		}(f)
	}
	// closer
	var total int64
	for size := range sizes {
		total += size
	}
	return total
}
