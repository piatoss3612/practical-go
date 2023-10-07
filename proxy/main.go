package main

import (
	"05-circuit-breaker/circuitbreaker"
	"context"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type proxyServer struct {
	cb *circuitbreaker.CircuitBreaker
	*http.Server
}

func (s *proxyServer) setup() {
	s.Server.Handler = s.handler()
}

func (s *proxyServer) handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := s.cb.Execute(r.Context(), func(ctx context.Context) (interface{}, error) {
			client := http.DefaultClient

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080", nil)
			if err != nil {
				return nil, err
			}

			return client.Do(req)
		})
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		httpResp := resp.(*http.Response)
		defer httpResp.Body.Close()

		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			slog.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(httpResp.StatusCode)
		_, _ = w.Write(body)
	})
}

func main() {
	cb := circuitbreaker.New(
		circuitbreaker.WithOpenTimeout(5*time.Second),
		circuitbreaker.WithStateChangeHook(
			func(from, to circuitbreaker.State) {
				slog.Info("state change", slog.String("from", from.String()), slog.String("to", to.String()))
			},
		),
		circuitbreaker.WithTripFunc(func(c circuitbreaker.Counter) bool {
			return c.TotalFailures > 3
		}),
	)

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
