c = 10
m = "Hello World!"

server:
	@echo "Starting TCP server..."
	go run ./listener

proxy-server:
	@echo "Starting TCP proxy server..."
	go run ./proxy

ping:
	@echo "Pinging TCP server..."
	go run ./client ping -c $(c)

echo:
	@echo "Sending echo request to TCP server..."
	go run ./client echo -c $(c) -m $(m)