# Format code 
fmt:
	go fmt ./...

# View possible issues in codebase
vet:
	go vet ./...

# Add any missing libraries and remove unsed ones
tidy: fmt
	go mod tidy

# Build the executable binary for the application
build:
	go build -o bin/

# Start the UDP Server
start: build
	./bin/dns-forward start

# Run the tests
test:
	go test -v ./cmd

# Clean project files and remove current binary in ./bin
clean:
	go clean
	rm ./bin/dns-forward


