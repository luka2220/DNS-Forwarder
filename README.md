# DNS Forwarder CLI
A CLI DNS Forwarder application that can resolve the IP address for a host either from its local cache or by forwarding the request to an authoritative nameserver.

## Testing the DNS Server
- Run the following commands below to ensure the DNS server is configured correctly.

- To test the DNS server
  * ` go run main.go start`
  * open a separate terminal window while the program is running to run the command below.
  * `dig @127.0.0.1 -p 8080 www.google.com`

- An example output:
  ![Screenshot 2024-01-16 at 10 00 23 AM](https://github.com/luka2220/DNS-Forwarder/assets/42144047/f3152265-f9be-4954-a5b2-06f6e1d47cf1)
