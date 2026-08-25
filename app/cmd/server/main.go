package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const homePage = `<!doctype html>
<html lang="es">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width,initial-scale=1">
	<title>Nodo comunitario</title>
	<style>
		:root {
			color-scheme: light;
			font-family: system-ui, sans-serif;
		}
		body {
			max-width: 52rem;
			margin: 0 auto;
			padding: 4rem 1.5rem;
			line-height: 1.6;
			color: #172033;
			background: #f4f7fb;
		}
		main {
			padding: 2rem;
			border: 1px solid #d7dfec;
			border-radius: .75rem;
			background: #fff;
		}
		.status {
			display: inline-block;
			padding: .25rem .65rem;
			border: 1px solid #9ab6df;
			border-radius: 999px;
			color: #174c8f;
			background: #edf5ff;
		}
	</style>
</head>
<body>
	<main>
		<p class="status">Nodo experimental</p>
		<h1>Nodo comunitario</h1>
		<p>
			Infraestructura autocustodiada para comunidades que necesitan
			publicar, comunicarse y preservar su memoria.
		</p>
		<ul>
			<li>Sin JavaScript.</li>
			<li>Sin rastreadores.</li>
			<li>Sin recursos externos.</li>
			<li>Accesible mediante Tor.</li>
		</ul>
	</main>
</body>
</html>`

var page = template.Must(template.New("home").Parse(homePage))

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'none'; style-src 'unsafe-inline'; "+
				"base-uri 'none'; form-action 'self'; frame-ancestors 'none'",
		)
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := page.Execute(w, nil); err != nil {
		log.Printf("template error: %v", err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprintln(w, "ok")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/healthz", health)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           securityHeaders(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("community onion app listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-stop

	log.Println("shutting down")
	_ = server.Close()
}
