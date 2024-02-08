package cmd

import (
	"log"
	"net"
	"os/exec"
	"sync"
	"testing"

	"github.com/miekg/dns"
)

func TestStartUDPServer(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Set up logger

	// Initialize a Mock UDP Server
	addr, _ := net.ResolveUDPAddr("udp", ":8080")
	conn, _ := net.ListenUDP("udp", addr)

	// Use a wait group to wait for the goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the UDP server in a goroutine
	go func() {
		// Start the UDP server to run asynchronously
		go startUDPServer(8080, conn, dns.DefaultServeMux)

		// Test the UDP server is running with dig command
		cmd := exec.Command("dig", "@127.0.0.1", "-p", "8080", "google.com")
		output, err := cmd.Output()

		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		log.Printf("\nOutput \n%s", output)

		// close the UDP connection
		conn.Close()

		wg.Done()
	}()

	// Wait for the goroutine to finish
	wg.Wait()
}
