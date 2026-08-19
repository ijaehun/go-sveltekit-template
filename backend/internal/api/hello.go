package api

import (
	"net/http"
	"time"
)

// GET /api/hello — 기본 JSON 응답 예제
func handleHello(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Hello from Go!",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// GET /api/hello/{name} — 경로 변수 예제 (Go 1.22+ PathValue)
func handleHelloName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Hello, " + name + "!",
	})
}
