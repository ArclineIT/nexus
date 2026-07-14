package handler

import "net/http"

// Ready handles GET /health and GET /ready.
func Ready() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, map[string]string{
			"status": "healthy",
		})
	}
}

