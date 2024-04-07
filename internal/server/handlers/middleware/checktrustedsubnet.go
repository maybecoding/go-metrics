package middleware

import (
	"net"
	"net/http"
)

func CheckTrustedSubnet(trustedSubNet *net.IPNet, ipAddrHeader string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Если требуется проверять доверенность узла, отправляющего метрику делаем это
			if trustedSubNet != nil {
				ipAddr := net.ParseIP(r.Header.Get(ipAddrHeader))
				if !trustedSubNet.Contains(ipAddr) {
					w.WriteHeader(http.StatusForbidden)
					return
				}
			}
			h.ServeHTTP(w, r)
		})
	}
}
