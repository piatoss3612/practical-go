package main

import "net/http"

func main() {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: h,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
