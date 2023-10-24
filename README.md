# 우아하게 종료하기 (Graceful Shutdown)

## 1. 로컬에서 실행하기

### 1.1. 실행하기

```bash
$ go run ./
```

### 1.2. 로그 확인하기

```bash
$ go run .
2023/10/24 17:36:11 DB connection established
2023/10/24 17:36:11 Starting server...
```

### 1.3. 서버 종료하기

```bash
$ go run .
2023/10/24 17:36:11 DB connection established
2023/10/24 17:36:11 Starting server...
^C2023/10/24 17:36:12 Shutting down server...
2023/10/24 17:36:12 DB connection closed
2023/10/24 17:36:12 Server shutdown complete
```

## 2. 컨테이너에서 실행하기

### 2.1. 컨테이너 빌드하기

```bash
$ docker build -t server .
```

### 2.2. 컨테이너 실행하기

```bash
$ docker run --name server -d -p 8080:8080 server
```

### 2.3. 컨테이너 중지하기

```bash
$ docker stop server
```

### 2.4. 로그 확인하기

```bash
$ docker logs server
2023/10/24 06:36:08 DB connection established
2023/10/24 06:36:08 Starting server...
2023/10/24 06:36:12 Shutting down server...
2023/10/24 06:36:12 DB connection closed
2023/10/24 06:36:12 Server shutdown complete
```