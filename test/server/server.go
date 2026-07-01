package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Tomrm18/rate-limiter/internal/clock"
	"github.com/Tomrm18/rate-limiter/pkg/algos"
)

func handler(message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(message))
	})
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", handler("Hello World\n"))

	bucketClock := clock.New()
	routeBucket := algos.NewBucketRateLimiter(1, 1, 1, 1*time.Second, bucketClock)

	log.Print("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", RateLimitMiddleware(routeBucket)(mux)))
}
