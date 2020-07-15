package main

import (
	"flag"
	"text/template"

	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var t *template.Template

func main() {

	// -verbose flag to set logging level to DebugLevel
	flagVerbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	if *flagVerbose {
		log.SetLevel(log.DebugLevel)
	}

	// Output to stdout instead of the default stderr
	// log.SetOutput(os.Stdout)

	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	r := mux.NewRouter()

	// Serve all static files in public directory
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// Define routes
	r.HandleFunc("/", parseTemplates(handlerIndex))

	// Set custom 404 page
	r.NotFoundHandler = http.HandlerFunc(handler404)

	// Set up the HTTP-server
	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func parseTemplates(h http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		// Parse all templates
		t, err = template.ParseGlob("./templates/*")
		if err != nil {
			log.Panic("Cannot parse templates", err)
		}

		h(w, r)
	}
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
	if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Warn(err)
	}
}

func handler404(w http.ResponseWriter, r *http.Request) {
	if err := t.ExecuteTemplate(w, "404.html", nil); err != nil {
		log.Warn(err)
	}
}
