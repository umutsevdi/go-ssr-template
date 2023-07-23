package api

import (
	"errors"
	"github.com/umutsevdi/webwatch/app/util"
	"mime"
	"net/http"
	"strings"
)

func Serve() {
	http.HandleFunc("/api/", index)
}

func index(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
	if ip == "" {
		ip = r.RemoteAddr
	}
	var file []byte
	var err error
	var status int = 200

	switch r.URL.Path {
	case "/api/":
		file = []byte("{\"value\": 5}")
	case "/api/data":
		err = errors.New("Not found")
	default:
		file = []byte("{\"path\":\"" + r.URL.Path + "\"}")
	}
	if err != nil {
		status = 404
	}
	contentHeader := mime.TypeByExtension(util.Ext(r.URL.Path))
	if contentHeader != "" {
		w.Header().Add("Content-Type", contentHeader)
	} else {
		w.Header().Add("Content-Type", "application/json")
	}
	w.WriteHeader(status)
	w.Write(file)
}
