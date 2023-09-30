# 잠자는 이발사 문제

## 0. 환경

- go 1.21.0

## 1. 실행

```bash
$ go run -race .
```

## 2. 테스트

```bash
$ go test -v -race -cover .
```

## 3. 테스트 결과

```bash
=== RUN   TestSleepingBarber
...
--- PASS: TestSleepingBarber (47.27s)
PASS
coverage: 100.0% of statements
ok      sleeping-barber 64.165s coverage: 100.0% of statements
```