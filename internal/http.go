package internal

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

// rawTemplate is the raw HTML template used for known packages
//
//go:embed resources/match.html
var rawTemplate string

// server is the HTTP server used by the application
var server *http.Server

// tmpl is the processed HTML template used for known packages
var tmpl *template.Template

func init() {
	// build the HTML template using the embedded file
	t, err := template.New("match").Parse(rawTemplate)
	if err != nil {
		slog.Error("failed to build HTML template", "error", err)
		panic(err)
	}
	tmpl = t
}

// Listen builds an HTTP server and starts listening for requests
func Listen(listen string) {
	// create a serve mux and register the handler
	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)

	// create the HTTP server itself using the serve mux
	server = &http.Server{
		Handler: mux,
		Addr:    listen,
	}

	// being listening for HTTP requests
	slog.Info("listening for HTTP requests", "listen", listen)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to listen for HTTP requests", "error", err)
			panic(err)
		}
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	to := "https://dsage.org"

	// check if the package is known to us
	if url, ok := packages[name]; ok {
		// ensure the client is ready to receive the response body
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		// build the template using the information we have
		err := tmpl.Execute(w, map[string]string{
			"ImportString": fmt.Sprintf("%s/%s git %s", r.Host, name, url),
			"URL":          url,
		})
		if err != nil {
			slog.Error("failed to execute HTML template", "error", err, "name", name)
			panic(err)
		}

		return
	}

	// redirect the client to the fallback URL
	w.Header().Set("Location", to)
	w.WriteHeader(http.StatusFound)
}

// Shutdown attempts to gracefully stop the HTTP server
func Shutdown() {
	// verify the server was started
	if server != nil {
		// trigger a shutdown with a 30-second timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				slog.Error("failed to stop HTTP server gracefully", "error", err)
				panic(err)
			}
		}
	}

	slog.Info("stopped HTTP server")
}
