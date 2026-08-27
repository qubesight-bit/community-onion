package main

import (
	"net/http"
)

const siteCSS = `
:root {
	color-scheme: light;
	--paper: #fbf8fa;
	--paper-raised: #ffffff;
	--ink: #1b1b1d;
	--muted: #44474c;
	--line: #75777d;
	--line-soft: #c5c6cd;
	--accent: #333e50;
	--accent-hover: #4a5568;
	--accent-text: #ffffff;
	--error: #93000a;
	--measure: 45rem;
	--radius: .25rem;
	font-family: Georgia, "Times New Roman", serif;
	font-size: 100%;
}

* {
	box-sizing: border-box;
}

html {
	background: var(--paper);
}

body {
	max-width: 72rem;
	margin: 0 auto;
	padding: 1.5rem;
	color: var(--ink);
	background: var(--paper);
	font-size: 1.125rem;
	line-height: 1.6;
}

main {
	max-width: var(--measure);
	margin: 2rem auto;
	padding: 2rem;
	border: 1px solid var(--line);
	border-radius: var(--radius);
	background: var(--paper-raised);
}

h1,
h2,
h3 {
	margin-top: 0;
	line-height: 1.2;
	letter-spacing: -.015em;
}

h1 {
	font-size: clamp(2rem, 7vw, 2.75rem);
}

h2 {
	font-size: clamp(1.4rem, 5vw, 1.8rem);
}

p,
ul {
	margin: 1rem 0;
}

a {
	color: var(--accent);
	text-decoration-thickness: .08em;
	text-underline-offset: .18em;
}

a:hover {
	color: var(--accent-hover);
}

a:focus-visible,
button:focus-visible,
input:focus-visible,
textarea:focus-visible,
select:focus-visible {
	outline: 3px solid var(--accent);
	outline-offset: 3px;
}

nav {
	padding-bottom: 1rem;
	border-bottom: 1px solid var(--line-soft);
}

.status {
	display: inline-block;
	margin: 0 0 1.25rem;
	padding: .3rem .7rem;
	border: 1px solid var(--accent);
	border-radius: var(--radius);
	color: var(--accent);
	font-size: .875rem;
	font-weight: 700;
	letter-spacing: .05em;
	text-transform: uppercase;
	background: transparent;
}

article {
	margin-top: 2rem;
	padding-top: 1.5rem;
	border-top: 1px solid var(--line);
}

article:first-of-type {
	margin-top: 1.5rem;
}

.content {
	white-space: pre-wrap;
	overflow-wrap: anywhere;
}

time {
	color: var(--muted);
	font-size: .9rem;
	font-weight: 700;
	letter-spacing: .03em;
}

button,
.button {
	display: inline-block;
	min-height: 44px;
	padding: .65rem 1rem;
	border: 2px solid var(--accent);
	border-radius: var(--radius);
	color: var(--accent-text);
	background: var(--accent);
	font: inherit;
	font-weight: 700;
	cursor: pointer;
}

button:hover,
.button:hover {
	color: var(--accent);
	background: transparent;
}

hr {
	margin: 2rem 0;
	border: 0;
	border-top: 1px solid var(--line);
}

::selection {
	color: var(--paper);
	background: var(--ink);
}

@media (max-width: 42rem) {
	body {
		padding: .75rem;
		font-size: 1rem;
	}

	main {
		margin: .75rem auto;
		padding: 1.25rem;
	}

	nav {
		line-height: 2;
	}
}

@media (prefers-reduced-motion: reduce) {
	* {
		scroll-behavior: auto !important;
		transition: none !important;
	}
}

@media print {
	body,
	main {
		max-width: none;
		margin: 0;
		padding: 0;
		border: 0;
		background: #fff;
	}

	nav,
	.status {
		display: none;
	}

	a {
		color: #000;
	}
}
`

func stylesheet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/assets/site.css" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")

	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	_, _ = w.Write([]byte(siteCSS))
}
