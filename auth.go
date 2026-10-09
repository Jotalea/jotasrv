package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
	"os"
)

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func newAuthHandler(handler http.Handler, user, pass string) http.Handler {
	userHash := sha256.Sum256([]byte(user))
	passHash := sha256.Sum256([]byte(pass))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if ok {
			uh := sha256.Sum256([]byte(u))
			ph := sha256.Sum256([]byte(p))
			userOK := subtle.ConstantTimeCompare(uh[:], userHash[:]) == 1
			passOK := subtle.ConstantTimeCompare(ph[:], passHash[:]) == 1
			if userOK && passOK {
				handler.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="jotasrv", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
