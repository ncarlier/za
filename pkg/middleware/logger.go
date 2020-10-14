package middleware

import (
	"net/http"
	"time"

	"github.com/ncarlier/za/pkg/logger"
)

type key int

const (
	requestIDKey key = 0
)

// Logger is a middleware to log HTTP request
func Logger(exceptions ...string) Middleware {
	except := make(map[string]bool)
	for _, path := range exceptions {
		except[path] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ignore := except[r.URL.Path]; !ignore {
				start := time.Now()
				defer func() {
					logger.Info.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), time.Since(start))
				}()
			}
			next.ServeHTTP(w, r)
		})
	}
}
