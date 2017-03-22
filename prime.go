package main

import (
    "fmt"
    "math"
)
/*
 * Prints out all prime numbers under max, then exits
 */
func main() {
    max := 10000000
    firstChannel := make(chan int)
    resultChannel := make(chan int)
    go worker(2, int(math.Sqrt(float64(max))), firstChannel, resultChannel)
    
    endChannel := make(chan int) // this is used as a signal for when we're done
    go printer(resultChannel, endChannel)  // need a printer goroutine so we don't deadlock
    for i := 2; i <= max; i++ {
        firstChannel <- i;
    }
    close(firstChannel) // tell the workers we're done
    <- endChannel // wait for the printer to signal it's done printing
}

/*
 * Reads values from result and prints them until it gets a value < 0, then it
 * signals end
 */
func printer(result <-chan int, end chan<- int){
    for i := range result {
        fmt.Println(i)
    }
    end <- 0
}

/*
 * One link in a long chain. Each link in the chain has base equal to a prime
 * number. Each link in the chain receives integers from myChannel and if that
 * number is not divisible by base, sends it to the next link in the chain. The
 * last link in the chain sends to the printer function, which prints it.
 */
func worker(base, end int, myChannel chan int, resultChannel chan int){
    var nextChannel chan int
    // we only need links up to sqrt(max)
    if base <= end {
        // make the new link
        nextChannel = make(chan int)
        go worker(<- myChannel, end, nextChannel, resultChannel)
    } else {
        // we're the last link, so send numbers that pass us to resultChannel
        nextChannel = resultChannel
    }
    // print our base
    resultChannel <- base
    
    // start receiving numbers to check
    for i := range myChannel{
        if i % base != 0 {
            // not divisible by base, so send it on to the next link
            nextChannel <- i
        }
    }
    close(nextChannel)
}