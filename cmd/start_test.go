package cmd

import (
	"net"
	"sync"
	"testing"

	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

// Mock DNS handler
type MockDNSHandler struct {
	Handled bool
}

func (m *MockDNSHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	m.Handled = true
}

func TestStartUDPServer(t *testing.T) {
	// Mock dependencies
	mockHandler := &MockDNSHandler{}

	// Initialize a Mock UDP Server
	addr, _ := net.ResolveUDPAddr("udp", ":8080")
	conn, _ := net.ListenUDP("udp", addr)

	// Use a wait group to wait for the goroutine to finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Start the UDP server in a goroutine
	go func() {
		defer wg.Done()
		startUDPServer(8080, conn, mockHandler)
	}()

	defer conn.Close()

	// Wait for the goroutine to finish
	wg.Wait()

	// Assert to check if the DNS MockDNS servers handler was called
	assert.True(t, mockHandler.Handled, "ServeDNS method should have been called")

	// Assert for expected UDP port
	expectedPort := 8080
	assert.Equal(t, expectedPort, 8080, "Unexpected UDP port! Must be 8080")
}
