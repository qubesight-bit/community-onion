package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"community-onion/internal/posts"

	"github.com/jackc/pgx/v5/pgxpool"
)

const postsPageHTML = `<!doctype html>
<html lang="es">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width,initial-scale=1">
	<title>Publicaciones — Nodo comunitario</title>
	<link rel="stylesheet" href="/assets/site.css">
</head>
<body>
	<main>
		<nav><a href="/">Inicio</a></nav>
		<h1>Publicaciones</h1>

		{{if .}}
			{{range .}}
			<article>
				<h2>{{.Title}}</h2>
				<p><time datetime="{{.PublishedAt.UTC.Format "2006-01-02T15:04:05Z"}}">{{.PublishedAtUTC}}</time></p>
				<div class="content">{{.Body}}</div>
			</article>
			{{end}}
		{{else}}
			<p>Todavía no hay publicaciones disponibles.</p>
		{{end}}
	</main>
</body>
</html>`

var postsPage = template.Must(
	template.New("published-posts").Parse(postsPageHTML),
)

func newPostsHandler(pool *pgxpool.Pool) http.HandlerFunc {
	repository := posts.NewRepository(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/publicaciones" {
			http.NotFound(w, r)
			return
		}

		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := timeContext(r, 5*time.Second)
		defer cancel()

		items, err := repository.ListPublished(ctx, 50)
		if err != nil {
			log.Printf("list published posts: %v", err)
			http.Error(w, "Servicio temporalmente no disponible", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if r.Method == http.MethodHead {
			w.WriteHeader(http.StatusOK)
			return
		}

		if err := postsPage.Execute(w, items); err != nil {
			log.Printf("posts template: %v", err)
		}
	}
}
