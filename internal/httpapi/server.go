package httpapi

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"homemaker-followup/internal/dashboard"
	"homemaker-followup/internal/domain"
	"homemaker-followup/internal/followup"
	"homemaker-followup/internal/reminders"
)

type Server struct {
	Service *followup.Service
	Clients []domain.Client
	Setting domain.ReminderSetting
	Today   time.Time
	mux     http.ServeMux
}

func NewServer(service *followup.Service, clients []domain.Client, setting domain.ReminderSetting, today time.Time) *Server {
	server := &Server{Service: service, Clients: clients, Setting: setting, Today: today}
	server.routes()
	return server
}

func (s *Server) Handler() http.Handler { return requestLogger(&s.mux) }

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.home)
	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/api/dashboard", s.dashboard)
	s.mux.HandleFunc("/api/records", s.records)
	s.mux.HandleFunc("/api/records/", s.recordDetail)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	template.Must(template.New("home").Parse(homeTemplate)).Execute(w, map[string]string{"Title": "家政客户回访簿"})
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	snapshot := dashboard.BuildSnapshot(s.Clients, s.Service.Records, s.Setting, s.Today)
	writeJSON(w, http.StatusOK, snapshot)
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		query := r.URL.Query().Get("q")
		writeJSON(w, http.StatusOK, s.Service.Search(query))
		return
	}
	if r.Method != http.MethodPost {
		allowMethods(w, "GET", "POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	form, err := decodeForm(r)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	record, err := form.ToRecord()
	if err != nil {
		writeError(w, statusForError(err), err)
		return
	}
	if err := s.Service.AddRecord(record, "web"); err != nil {
		writeError(w, statusForError(err), err)
		return
	}
	writeJSON(w, http.StatusCreated, record)
}

func (s *Server) recordDetail(w http.ResponseWriter, r *http.Request) {
	id := parseID(r.URL.Path)
	for _, record := range s.Service.Records {
		if record.ID == id {
			writeJSON(w, http.StatusOK, record)
			return
		}
	}
	http.NotFound(w, r)
}

func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

var homeTemplate = `<!doctype html><html lang="zh-CN"><head><meta charset="utf-8"><title>{{.Title}}</title></head><body><main><h1>{{.Title}}</h1><p>请打开 web/index.html 查看提醒看板。</p></main></body></html>`

func reminderSummary(records []domain.FollowUpRecord, setting domain.ReminderSetting, today time.Time) []reminders.Reminder {
	return reminders.Build(records, setting, today)
}
