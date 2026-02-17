package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"oneproxy-clientwebui/internal/config"
)

func NewProxyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cfg := config.Get()
		target, err := url.Parse(strings.TrimRight(cfg.APIBaseURL, "/"))
		if err != nil {
			http.Error(w, "invalid upstream URL", http.StatusBadGateway)
			return
		}

		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = target.Scheme
				req.URL.Host = target.Host
				req.Host = target.Host
			},
			FlushInterval: -1,
			ModifyResponse: func(resp *http.Response) error {
				resp.Header.Del("Access-Control-Allow-Origin")
				resp.Header.Del("Access-Control-Allow-Methods")
				resp.Header.Del("Access-Control-Allow-Headers")
				return nil
			},
		}

		proxy.ServeHTTP(w, r)
	}
}
