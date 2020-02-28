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

		type Content struct {
			Title       string
			Description string
			Datetime    string
			Datetext    string
		}

		type Topic struct {
			Image   string
			Alt     string
			Content Content
		}

		type Layout1 struct {
			GlobalNav []string
			HotTopic  Topic
		}

		l1 := Layout1{
			GlobalNav: []string{"HOME", "ABOUT", "NEWS", "TOPICS", "DOCS", "BLOG"},
			HotTopic: Topic{
				Image: "assets/elephant-4389434_1280.jpg", Alt: "elephant",
				Content: Content{
					Title:       "The water is literally pouring",
					Description: "Elephant is just the tip of the iceberg, but considering that it remains hard to point to an unequivocal link between climate change and more frequent severe weather events.",
					Datetime:    "2020-02-27",
					Datetext:    "2020.02.27 TUE",
				}},
		}

		t := template.Must(template.ParseFiles(filepath.Join(templatesDir, "layout1.html")))
		_ = t.ExecuteTemplate(w, "layout1.html", l1)
	})

	http.HandleFunc("/layout2", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		type Layout2 struct {
			Title       string
			Subtitle    string
			Description string
		}

		l2 := Layout2{
			Title:       "NeQUE PORRO QUISQUAM",
			Subtitle:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit",
			Description: "sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		}

		t := template.Must(template.ParseFiles(filepath.Join(templatesDir, "layout2.html")))
		_ = t.ExecuteTemplate(w, "layout2.html", l2)
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
