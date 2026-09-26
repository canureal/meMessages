package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                now := time.Now()
                requesturi := r.RequestURI
                method := r.Method

                log.Printf("%v (%v), at %v", requesturi, method, now)
                next.ServeHTTP(w, r)
        })
}
