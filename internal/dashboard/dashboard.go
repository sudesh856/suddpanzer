package dashboard

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
	_ "embed"

	"github.com/gorilla/websocket"
)

//go:embed dashboard.html
var DashboardHTML []byte
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type MetricsSnapshot struct {
	Timestamp  int64   `json:"timestamp"`
	RPS        float64 `json:"rps"`
	P50        int64   `json:"p50"`
	P95        int64   `json:"p95"`
	P99        int64   `json:"p99"`
	ErrorRate  float64 `json:"error_rate"`
	TotalReqs  int64   `json:"total_requests"`
	ActiveVUs  int     `json:"active_vus"`
}

type Server struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]bool
	cancel  func()
}

func New(cancel func()) *Server {
	return &Server{
		clients: make(map[*websocket.Conn]bool),
		cancel:  cancel,
	}
}

func (s *Server) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
	}()

	// keep connection alive
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (s *Server) HandleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	s.cancel()
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Broadcast(snap MetricsSnapshot) {
	data, err := json.Marshal(snap)
	if err != nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for conn := range s.clients {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (s *Server) Start(addr string, htmlContent []byte) {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", s.HandleWS)
	mux.HandleFunc("/api/stop", s.HandleStop)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(htmlContent)
	})
	srv := &http.Server{Addr: addr, Handler: mux}
	go srv.ListenAndServe()
}

func (s *Server) StartBroadcasting(getSnapshot func() MetricsSnapshot, ctx interface{ Done() <-chan struct{} }) {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.Broadcast(getSnapshot())
			}
		}
	}()
}