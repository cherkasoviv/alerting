package hash

import (
	"net/http"
)

func Hash256Middleware(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

		}

		return http.HandlerFunc(fn)
	}
}
