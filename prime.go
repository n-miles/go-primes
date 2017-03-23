package main

import (
	"fmt"
	"math"
	"sync"
	"time"
)

var wg sync.WaitGroup

/*
 * Prints out all prime numbers under max, then exits
 */
func main() {
	startTime := time.Now()
	max := 10000000
	firstChannel := make(chan int, 100)
	resultChannel := make(chan int, 100)
	// add the first worker (assume we know 2 is prime)
	go worker(2, int(math.Sqrt(float64(max))), firstChannel, resultChannel)
	// workers feed the printer
	go printer(resultChannel)
	wg.Add(1)
	for i := 2; i <= max; i++ {
		firstChannel <- i
	}
	// tell the workers we're done
	close(firstChannel)
	// wait for the printer to finish printing all the numbers
	wg.Wait()
	timeElapsed := time.Now().Sub(startTime)
	fmt.Println(timeElapsed)
}

/*
 * Reads values from result and prints them until it gets a value < 0, then
 * signals main() and exits
 */
func printer(result <-chan int) {
	for i := range result {
		fmt.Println(i)
	}
	wg.Done()
}

/*
 * One link in a long chain. Each link in the chain has base equal to a prime
 * number. Each link in the chain receives integers from myChannel and if that
 * number is not divisible by base, sends it to the next link in the chain. The
 * last link in the chain sends to the printer function, which prints it. When
 * a link's channel is closed, it then closes the next channel in the chain. The
 * last one closes printer()'s channel, which triggers the program exit.
 */
func worker(base, end int, myChannel <-chan int, resultChannel chan int) {
	var nextChannel chan int
	// we only need links up to sqrt(max)
	if base <= end {
		// make the new link
		nextChannel = make(chan int, 100)
		go worker(<-myChannel, end, nextChannel, resultChannel)
	} else {
		// we're the last link, so send numbers that pass us to resultChannel
		nextChannel = resultChannel
	}
	// print our base
	resultChannel <- base

	// start receiving numbers to check
	for i := range myChannel {
		if i%base != 0 {
			// not divisible by base, so send it on to the next link
			nextChannel <- i
		}
	}
	close(nextChannel)
}
