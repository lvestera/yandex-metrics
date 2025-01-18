package sign

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
)

var keyStr string

func NewSign(key string) {
	keyStr = key
}

func RequestHashCheck(h http.Handler) http.Handler {
	checkHashFn := func(w http.ResponseWriter, r *http.Request) {

		hash := r.Header.Get("HashSHA256")

		if len(hash) > 0 {

			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			defer r.Body.Close()
			r.Body = io.NopCloser(bytes.NewBuffer(body))

			requestHash, err := hex.DecodeString(hash)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			newHash := CalcHash(body, keyStr)
			if !hmac.Equal(requestHash, newHash) {
				http.Error(w, "incorrect hash", http.StatusBadRequest)
			}
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(checkHashFn)
}

func CalcHash(body []byte, key string) []byte {
	/*
	   Реализуйте механизм подписи передаваемых данных по алгоритму SHA256. Для этого посчитайте hash от всего тела запроса и разместите его в HTTP-заголовке HashSHA256.
	   Хеш нужно считать от строки с учётом ключа, который передан агенту/серверу на старте: hash(value, key)
	*/

	h := hmac.New(sha256.New, []byte(key))
	h.Write(body)
	return h.Sum(nil)
}
