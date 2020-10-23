package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/gorilla/mux"
)

// Server is the Vote server.
type Server struct {
	Router *mux.Router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.HandleFunc("/_healthcheck", s.handleHealthCheck())
	s.Router.HandleFunc("/", s.handleRoot()).Methods(http.MethodGet)

	s.Router.ServeHTTP(w, r)
}

func (s *Server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emojis, err := retrieveEmojis()
		if err != nil {
			emojis = nil
		}
		sorter := &emojiSorter{
			emojis: emojis,
			by: func(e1, e2 emoji) bool {
				return e1.Count < e2.Count
			},
		}
		sort.Sort(sorter)
		log.Printf("emojis: %v\n", emojis)
		renderTemplate(w, "index", emojis)
	}
}

type emoji struct {
	Emoji string `json:"emoji"`
	Count int `json:"count"`
}

type emojiSorter struct {
	emojis []emoji
	by func(e1, e2 emoji) bool
}

func (s *emojiSorter) Len() int {
	return len(s.emojis)
}

func (s *emojiSorter) Swap(i, j int) {
	s.emojis[i], s.emojis[j] = s.emojis[j], s.emojis[i]
}

func (s *emojiSorter)  Less(i, j int) bool {
	return s.by(s.emojis[i], s.emojis[j])
}

func retrieveEmojis() ([]emoji, error) {
	endpoint := fmt.Sprintf("http://api.%s:8080/tweets/emojis/", os.Getenv("COPILOT_SERVICE_DISCOVERY_ENDPOINT"))
	resp, err := http.Get(endpoint)
	if err != nil {
		log.Printf("WARN: server: coudln't get emojis: %v\n", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("WARN: server: get vote response status: %d\n", resp.StatusCode)
		return nil, errors.New("unexpected status code")
	}
	defer resp.Body.Close()
	data := struct {
		Emojis []emoji `json:"emojis"`
	}{}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&data); err != nil {
		log.Printf("ERROR: server: decode emojis: %v\n", err)
		return nil, fmt.Errorf("server: decode emojis: %v",err)
	}
	log.Printf("INFO: server: received %d emojis\n", len(data.Emojis))
	return data.Emojis, nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, emojis []emoji) {
	f := filepath.Join("templates", tmpl + ".html")
	t, err := template.New(path.Base(f)).
		Funcs(template.FuncMap{
			"Iterate": func(n int) []int {
				var arr []int
				for i := 1; i <= n; i += 1 {
					arr = append(arr, i)
				}
				return arr
			},
			"Sum": func(emojis []emoji) int {
				var sum int
				for _, emoji := range emojis {
					sum += emoji.Count
				}
				return sum
			},
			"Percentage": func(n, total int) int {
				return n*100/total
			},
		}).
		ParseFiles(f)
	if err != nil {
		log.Fatalf("parse file: %v\n", err)
	}
	if err := t.Execute(w, struct {
		Emojis []emoji
	} {
		Emojis: emojis,
	}); err != nil {
		log.Fatalf("execute template: %v\n", err)
	}
}