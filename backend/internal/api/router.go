// Package api: 라우트 등록. feature별로 파일을 나누고 여기서 mount한다.
package api

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	// Go 1.22+ 패턴 라우팅: "METHOD /path" 형식, {id} 같은 경로 변수 지원
	mux.HandleFunc("GET /api/health", handleHealth)
	mux.HandleFunc("GET /api/hello", handleHello)
	mux.HandleFunc("GET /api/hello/{name}", handleHelloName)
	mux.HandleFunc("GET /api/events", handleEvents) // SSE 데모

	// 새 feature 추가 시: internal/api/<feature>.go 만들고 여기에 등록
	// mux.HandleFunc("POST /api/todos", handleCreateTodo)
}
