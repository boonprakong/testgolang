package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

// Define constants
const (
	args                  = 2  // Expected number of command line arguments
	seconds time.Duration = 10 // Timeout duration for TCP connection
)

// Function for testing TCP connection to a given IP address and port
func testTCPConnection(ip string, port int, doneChannel chan bool) {
	// Try to establish TCP connection to given IP address and port with a timeout
	_, err := net.DialTimeout(
		"tcp",
		ip+":"+strconv.Itoa(port),
		time.Second*seconds,
	)

	// If connection is successful, print message indicating that the port is open
	if err == nil {
		fmt.Printf("Port %d: Open\n", port)
	}

	doneChannel <- true // Signal that goroutine is done by writing to channel
}

func main() {
	// Check if the number of command line arguments is correct
	if len(os.Args) != args {
		// Print usage message and exit with error code
		fmt.Fprintf(os.Stderr, "Using: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}

	// Get the IP address from command line arguments
	target := os.Args[1]

	// Initialize counter for active threads and channel for signaling when all threads are done
	activeThreads := 0
	doneChannel := make(chan bool)

	// Launch a goroutine for testing TCP connection for each port (0 to 65535)
	for port := 0; port <= 65535; port++ {
		go testTCPConnection(target, port, doneChannel) // Launch goroutine for testing TCP connection
		activeThreads++                                 // Increment counter for active threads
	}

	// Wait for all goroutines to finish
	for activeThreads > 0 {
		<-doneChannel   // Read from channel to wait for a goroutine to finish
		activeThreads-- // Decrement counter for active threads
	}
}
