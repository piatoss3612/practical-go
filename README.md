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
--- PASS: TestSleepingBarber (44.51s)
PASS
coverage: 98.4% of statements
ok      sleeping-barber 61.099s coverage: 98.4% of statements
```