package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
)

func Hash256Middleware(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ow := w

			if key != "" {
				requestHashSHA256 := r.Header.Get("HashSHA256")
				fmt.Println(requestHashSHA256)
				fmt.Println(key)
				b, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				h := hmac.New(sha256.New, []byte(key))
				h.Write(b)
				checkResult := h.Sum(nil)
				if !hmac.Equal([]byte(requestHashSHA256), checkResult) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

			}
			next.ServeHTTP(ow, r)
		}

		return http.HandlerFunc(fn)
	}
}
