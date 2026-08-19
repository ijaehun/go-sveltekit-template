package api

import (
	"fmt"
	"net/http"
	"time"
)

// GET /api/events — Server-Sent Events(SSE) 데모.
// Go의 실시간 스트리밍 손맛을 보여주는 예제: 1초마다 서버 시간을 push한다.
// 프론트에서는 EventSource로 구독한다 (routes/+page.svelte 참고).
// 웹소켓이 필요해지면 coder/websocket 패키지 추가를 고려.
func handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming unsupported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done(): // 클라이언트가 연결을 끊으면 종료
			return
		case t := <-ticker.C:
			fmt.Fprintf(w, "data: %s\n\n", t.Format("15:04:05"))
			flusher.Flush()
		}
	}
}
