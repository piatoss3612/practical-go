# 01 tcp

## 1. Run the server

```bash
$ go run ./server
```

or

```bash
$ make server
```

## 2. Run the proxy server (optional)

```bash
$ go run ./proxy
```

or

```bash
$ make proxy-server
```

## 3. Run the client

### 3.1 Ping

```bash
$ go run ./client ping -c <count>
```

or

```bash
$ make ping c=<count>
```

### 3.2 Echo the message

```bash
$ go run ./client echo -m <message> -c <count>
```

or

```bash
$ make echo m=<message> c=<count>
```

### 3.3 Proxy (Proxy server must be running)

```bash
$ go run ./client proxy -cmd <command> -m <message> -c <count>
```

or

```bash
$ make proxy-client cmd=<command> m=<message> c=<count>
```