package pages

import (
	"log"
	"mime"
	"net/http"
	"strings"

	"github.com/umutsevdi/webwatch/app/util"
)

// Contains ready to render page templates, after being merged with components
var pageCaches map[string][]byte

func Serve() {
	util.StartIndexing()
	util.MapEach("pages", func(k string, v util.FData) {
		http.HandleFunc(k, ServePage)
	})
	util.MapEach("static", func(k string, v util.FData) {
		http.HandleFunc(k, ServeStatic)
	})
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
	if ip == "" {
		ip = r.RemoteAddr
	}
	var file []byte
	var status int = 200

	if fData := util.MapGet("pages", r.URL.Path); fData != nil {
		file = fData.Content
	} else {
		switch r.URL.Path {
		case "/robots.txt":
			if fData := util.MapGet("static", "/static/robots.txt"); fData != nil {
				file = fData.Content
			}
		case "/favicon.ico":
			if fData := util.MapGet("static", "/static/favicon.ico"); fData != nil {
				file = fData.Content
			}
		case "/sitemap.xml":
			file = []byte(sitemap())
		}
	}
	if file == nil {
		if fData := util.MapGet("pages", "/not-found"); fData != nil {
			file = fData.Content
		}
		status = 404
	}
	contentHeader := mime.TypeByExtension(util.Ext(r.URL.Path))
	if contentHeader != "" {
		w.Header().Add("Content-Type", contentHeader)
	} else {
		w.Header().Add("Content-Type", "text/html")
	}
	log.Println("GET:", ip, r.URL.Path, status)
	w.WriteHeader(status)
	w.Write(file)
}

func ServeStatic(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
	if ip == "" {
		ip = r.RemoteAddr
	}
	var file []byte
	var status int = 200

	if fData := util.MapGet("static", r.URL.Path); fData != nil {
		file = fData.Content
	} else {
		status = 404
		log.Println("GET:", ip, r.URL.Path, status)
	}
	// Get mime type value and insert it into content header if it exists
	if h := mime.TypeByExtension(util.Ext(r.URL.Path)); h != "" {
		w.Header().Add("Content-Type", h)
	}
	w.WriteHeader(status)
	w.Header().Set("Cache-Control", "max-age=3600")
	w.Write(file)
}

func sitemap() string {
	var s strings.Builder = strings.Builder{}
	s.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?> <urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"><url><loc>")
	s.WriteString(*util.C.URI + "/</loc><priority>1.0</priority></url>")

	util.MapEach(
		"pages",
		func(p string, v util.FData) {
			if len(p) > 0 && p != "/not-found.html" && p != "/" {
				s.WriteString("<url><loc>" + *util.C.URI + p + "</loc> <priority>0.8</priority> </url>")
			}
		})

	s.WriteString("</urlset>")
	return s.String()
}
