package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
			fmt.Println("Goroutine finish.")
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("Main routinue finish")
}
