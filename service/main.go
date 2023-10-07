package main

import "net/http"

func main() {
	// '/' 경로로 요청이 들어오면 "Hello, world!"를 응답으로 보내는 핸들러를 생성한다.
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// 8080 포트와 h를 핸들러로 갖는 http.Server를 생성한다.
	srv := http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	// http.Server를 실행한다.
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
