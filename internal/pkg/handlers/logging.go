package handlers

import (
	"net/http"
	"strings"

	"github.com/felixge/httpsnoop"
	log "github.com/sirupsen/logrus"
)

type LoggingHandler struct {
	Handler http.Handler
}

func (h LoggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m := httpsnoop.CaptureMetrics(h.Handler, w, r)
	l := log.WithFields(log.Fields{
		"status": m.Code,
		"size":   m.Written,
		"took":   m.Duration,
		"method": r.Method,
		"url":    r.RequestURI,
		"remote": r.RemoteAddr,
	})
	if strings.HasPrefix(r.RequestURI, "/static") {
		return
	}

	if m.Code > 400 && m.Code < 500 {
		l.Level = log.WarnLevel
	} else if m.Code >= 500 {
		l.Level = log.ErrorLevel
	} else {
		l.Level = log.InfoLevel
	}

	l.Print("request handled")
}
