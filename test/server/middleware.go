package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/Tomrm18/rate-limiter/internal/limiter"
)

func RateLimitMiddleware(bucket limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Print("request received in rate limit middleware")

			// check the request is allowed
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ip = r.RemoteAddr
			}
			fmt.Println("ip = ", ip)

			res, err := bucket.Allow(ip)
			if err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if res.Success {
				next.ServeHTTP(w, r)
			} else {
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(int(res.Limit)))
				w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(int(res.TokensRemaining)))
				w.Header().Set("X-RateLimit-Reset", strconv.Itoa(int(res.TimeUntilRefill.Seconds())))
				log.Print("rejecting request with 429")
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
		})
	}
}
