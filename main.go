package main

import (
	"compress/gzip"
	"flag"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
)

const assetsDir = "assets"
const templatesDir = "templates"

func main() {
	var addr = flag.String("addr", ":5963", "The address of the application")
	flag.Parse()

	http.HandleFunc("/layout1", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		t := template.Must(template.ParseFiles(filepath.Join(templatesDir, "layout1.html")))
		_ = t.ExecuteTemplate(w, "layout1.html", nil)
	})

	http.Handle("/assets/", gzippify(http.StripPrefix("/assets/", http.FileServer(http.Dir(assetsDir)))))

	log.Println("Server listening on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzippify(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gzw := gzip.NewWriter(w)
		defer gzw.Close()
		grw := gzipResponseWriter{Writer: gzw, ResponseWriter: w}
		h.ServeHTTP(grw, r)
	})
}
