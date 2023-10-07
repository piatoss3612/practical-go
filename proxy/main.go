package main

import (
	"05-circuit-breaker/circuitbreaker"
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// proxyServer는 서킷 브레이커를 사용하는 프록시 서버를 나타내는 구조체
type proxyServer struct {
	cb           *circuitbreaker.CircuitBreaker // 서킷 브레이커
	*http.Server                                // 프록시 서버
}

// setup은 프록시 서버를 설정하는 메서드
func (s *proxyServer) setup() {
	s.Server.Handler = s.handler()
}

// handler는 프록시 서버의 핸들러를 반환하는 메서드
func (s *proxyServer) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 서킷 브레이커로 wrapping될 함수
		fn := func(ctx context.Context) (interface{}, error) {
			client := http.DefaultClient // http.Client 생성

			// http.Client를 사용하여 http://localhost:8080 경로로 GET 요청을 보낸다.
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
			if err != nil {
				return nil, err
			}

			return client.Do(req) // http.Response를 반환한다.
		}

		// 서킷 브레이커로 wrapping된 함수를 실행한다.
		resp, err := s.cb.Execute(r.Context(), fn)
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		httpResp := resp.(*http.Response) // 타입 assertion을 사용하여 http.Response를 얻는다.
		defer httpResp.Body.Close()

		w.WriteHeader(httpResp.StatusCode) // http.Response의 status code를 응답한다.
		_, _ = io.Copy(w, httpResp.Body)   // http.Response의 body를 응답한다.
	})
}

func main() {
	// 서킷 브레이커 생성
	cb := circuitbreaker.New(
		circuitbreaker.WithOpenTimeout(5*time.Second), // open 상태에서 half open 상태로 전환되기 위한 시간을 5초로 설정
		circuitbreaker.WithStateChangeHook(
			func(from, to circuitbreaker.State) {
				slog.Info("state change", slog.String("from", from.String()), slog.String("to", to.String()))
			},
		), // 서킷 브레이커의 상태가 변경될 때 호출되는 함수를 설정
		circuitbreaker.WithTripFunc(func(c circuitbreaker.Counter) bool {
			return c.TotalFailures > 3
		}), // 서킷 브레이커가 open 상태로 전환되기 위한 조건을 판단하는 함수를 설정
	)

	// 프록시 서버 생성 및 실행
	srv := &proxyServer{
		cb: cb,
		Server: &http.Server{
			Addr: ":4000",
		},
	}

	srv.setup()

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
