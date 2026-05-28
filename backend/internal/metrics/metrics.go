package metrics

import (
	"fmt"
	"net/http"
	"sync"
)

type Registry struct {
	mu       sync.RWMutex
	counters map[string]int64
}

func NewRegistry() *Registry {
	return &Registry{counters: map[string]int64{}}
}

func (r *Registry) Inc(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counters[name]++
}

func (r *Registry) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		r.mu.RLock()
		defer r.mu.RUnlock()
		for name, value := range r.counters {
			fmt.Fprintf(w, "findmesh_%s_total %d\n", name, value)
		}
	})
}
