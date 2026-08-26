package main

import (
	"context"
	"net/http"
	"time"
)

func timeContext(
	r *http.Request,
	timeout time.Duration,
) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), timeout)
}
