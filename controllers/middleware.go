package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/leesper/holmes"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

func LoggerMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rw := NewResponseWriter(w)
	start := time.Now()
	logging := fmt.Sprintf("%s -- %v %s %s %s %s - %s %v",
		r.RemoteAddr,
		start,
		r.Method,
		r.URL.Path,
		r.Proto,
		http.StatusText(rw.StatusCode),
		r.Header.Get("User-Agent"),
		time.Since(start))

	holmes.Infoln(logging)

	next(rw, r)
}

//CorsMiddleware handle CORS request, see https://github.com/rs/cors
func CorsMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// holmes.Infoln(r.Header)
	if "OPTIONS" == r.Method {
		headers := w.Header()
		origin := r.Header.Get("Origin")
		// Always set Vary headers
		// see https://github.com/rs/cors/issues/10,
		//     https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
		headers.Add("Vary", "Origin")
		headers.Add("Vary", "Access-Control-Request-Method")
		headers.Add("Vary", "Access-Control-Request-Headers")

		reqMethod := r.Header.Get("Access-Control-Request-Method")
		reqHeaders := strings.Split(r.Header.Get("Access-Control-Request-Headers"), ",")
		headers.Set("Access-Control-Allow-Origin", origin) // 信任所有来源
		// Spec says: Since the list of methods can be unbounded, simply returning the method indicated
		// by Access-Control-Request-Method (if supported) can be enough
		headers.Set("Access-Control-Allow-Methods", strings.ToUpper(reqMethod)) // 允许请求的方法
		if len(reqHeaders) > 0 {
			// Spec says: Since the list of headers can be unbounded, simply returning supported headers
			// from Access-Control-Request-Headers can be enough
			headers.Set("Access-Control-Allow-Headers", strings.Join(reqHeaders, ", ")) // 允许请求的自定义头参数
		}
		w.WriteHeader(http.StatusOK)
	} else {
		next(w, r)
	}
}
