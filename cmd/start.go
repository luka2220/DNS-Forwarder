/*
Copyright © 2024 Luka Piplica piplicaluka64@gmail.com
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a UDP server.",
	Long: `The start commands starts a UDP server.
		By default the server starts on port 8080.
		i.e
		dns-forwarder start
	`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetFlags(log.LstdFlags | log.Lshortfile) // Set up logger

		fmt.Println("Starting UDP server...")
		var port int = 8080

		// Creating a UDP address to listen on all available network interfaces
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
		if err != nil {
			log.Fatal(err)
		}

		// Creating a UDP listener
		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			log.Fatal(err)
		}

		// Call the function responsible for starting the UDP server
		startUDPServer(port, conn, dns.DefaultServeMux)
	},
}

func startUDPServer(port int, conn *net.UDPConn, handler dns.Handler) {
	dns.HandleFunc(".", handleDNSRequest)

	// Close the UDP server once Run is finished executing
	defer conn.Close()

	fmt.Printf("UDP server is listening on port %d\n", port)

	// Creating a tcp listener
	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}

	 dns.ActivateAndServe(l, conn, handler)
}

var cacheMutex sync.RWMutex
var dnsCache = make(map[string]dns.Msg)

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	fmt.Println("DNS Message: ", r)

	question := r.Question[0]

	cacheKey := fmt.Sprintf("%s|%d", question.Name, question.Qtype)

	// Check the cache for a previous response
	cacheMutex.RLock()
	cachedResponse, cacheExists := dnsCache[cacheKey]
	cacheMutex.RUnlock()

	if cacheExists {
		// Serve the response from the cache
		m.Answer = append(m.Answer, cachedResponse.Answer...)
	} else {
		// Forward DNS query to another DNS server
		forwardedResponse, err := dns.Exchange(r, "8.8.8.8:53")
		if err != nil {
			fmt.Println("Error forwarding DNS request:", err)
			m.SetRcode(r, dns.RcodeServerFailure)
			w.WriteMsg(m)
			return
		}

		// Check the query type (A record or IPv4 address)
		if question.Qtype == dns.TypeA {
			// Serve the response from the cache
			m.Answer = append(m.Answer, cachedResponse.Answer...)

			// Cache the response
			cacheMutex.Lock()
			dnsCache[cacheKey] = *forwardedResponse
			cacheMutex.Unlock()
		} else {
			// Handle other query types or respond with an error
			m.SetRcode(r, dns.RcodeNameError)
		}

		// Print that the request has been forwarded
		fmt.Println("Forwarded DNS request to 8.8.8.8:53")
	}

	w.WriteMsg(m)
}

func init() {
	rootCmd.AddCommand(startCmd)
}
