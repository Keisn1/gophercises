// package main

// import (
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"sync"
// 	"syscall"
// 	"time"
// )

// func userInput(stop chan struct{}, wg *sync.WaitGroup) {
// 	var input string
// 	fmt.Println("Waiting for user input (type something and press enter):")

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		for {
// 			scanned, err := fmt.Scanln(&input)
// 			if err != nil {
// 				fmt.Println("Error occurred while reading input:", err)
// 			}
// 			fmt.Println(scanned)
// 		}
// 	}()

// 	select {
// 	case <-stop:
// 		fmt.Println("Function stopped")
// 		return
// 	}
// }

// func main() {
// 	stop := make(chan struct{}) // Channel to signal stop
// 	var wg sync.WaitGroup       // WaitGroup for synchronization

// 	// Setting up signal handling for interruption (e.g., Ctrl+C)
// 	sig := make(chan os.Signal, 1)
// 	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

// 	go userInput(stop, &wg)

// 	timeout := time.After(10 * time.Second) // Set timeout duration

// 	select {
// 	case <-timeout:
// 		fmt.Println("Timeout occurred after 30 seconds. Stopping input reading.")
// 		close(stop) // Signal the userInput Goroutine to stop
// 	case <-sig:
// 		fmt.Println("Received interrupt signal, stopping the function immediately")
// 		close(stop) // Signal the userInput Goroutine to stop
// 	}

// 	wg.Wait() // Wait for userInput Goroutine to complete

// 	// Allow some time for the function to stop
// 	time.Sleep(1 * time.Second)
// }
