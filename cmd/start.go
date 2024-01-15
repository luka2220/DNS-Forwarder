/*
Copyright © 2024 Luka Piplica piplicaluka64@gmail.com
*/
package cmd

import (
	"fmt"
	"net"

	"program/dns-forward/error"

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
		fmt.Println("Starting UDP server...")
		var port int = 8080

		// Creating a UDP address to listen on all available network interfaces
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
		error.Check(err)

		// Creating a UDP listener
		conn, err := net.ListenUDP("udp", addr)
		error.Check(err)
		// Close the UDP server once Run in finished executing
		defer conn.Close()

		fmt.Printf("UDP server is listening on port %d\n", port)

		// Buffer to hold incoming data
		buffer := make([]byte, 1024)

		for {
			// Read data from the connection
			n, clientAddr, err := conn.ReadFromUDP(buffer)
			error.Check(err)

			fmt.Printf("Received %d bytes from %s: %s\n", n, clientAddr, buffer[:n])

			// Respond to the client
			response := []byte("Hello from UDP server!")
			_, err = conn.WriteToUDP(response, clientAddr)
			error.Check(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
