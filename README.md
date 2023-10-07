# Circuit Breaker 구현

## 0. 환경

- go 1.21.0
- postman (버전은 상관 없음)

## 1. 실행

### service 실행

```bash
$ go run ./service
```

### proxy 실행

- `proxy/main.go` 파일의 `main` 함수에서 호출된 `curcuitbreaker.New` 함수의 옵션은 기호에 따라 변경하여 실행 가능

```bash
$ go run ./proxy
```

### 요청 보내기

1. `postman`을 이용하여 `proxy` 서버에 요청을 보낸다. (포트는 `4000`, 경로는 `/`)
2. `proxy` 서버는 `service` 서버에 요청을 보내고, `service` 서버는 `proxy` 서버에 응답을 보낸다.
3. `proxy` 서버는 `service` 서버로부터 받은 응답을 `postman`에게 전달한다.
4. `service` 서버를 종료하고 `postman`을 통해 다시 요청을 보내면 `500 Internal Server Error`를 반환한다.
5. 요청이 여러 번 실패하여 서킷 브레이커의 `trip` 함수에 지정된 threshold를 넘으면 서킷 브레이커는 `open` 상태가 되고, `service` 서버로의 요청을 거부한다.
6. `open` 상태가 된 서킷 브레이커는 `timeout` 시간이 지나면 `half-open` 상태가 되고, `service` 서버로의 요청을 일시적으로 허용한다.
7. `half-open` 상태가 된 서킷 브레이커는 `halfOpenMaxSuccesses` 횟수만큼 `service` 서버로의 요청이 성공하면 `close` 상태가 되고, `service` 서버로의 요청을 허용한다.