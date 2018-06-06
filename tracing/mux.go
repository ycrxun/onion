package tracing

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing/opentracing-go"
)

// TracedServeMux is a wrapper around mux.Router that instruments handlers for tracing.
type TracedServeMux struct {
	mux    *mux.Router
	tracer opentracing.Tracer
}

// NewServeMux creates a new TracedServeMux.
func NewServeMux(tracer opentracing.Tracer) *TracedServeMux {
	return &TracedServeMux{
		mux:    mux.NewRouter(),
		tracer: tracer,
	}
}

// Handle implements mux.Router#Handle
func (tm *TracedServeMux) Handle(pattern string, handler http.Handler) {
	middleware := nethttp.Middleware(
		tm.tracer,
		handler,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + " " + r.URL.Path
		}),
	)
	tm.mux.Handle(pattern, middleware)
}

// HandleFunc implements mux.Router#HandleFunc
func (tm *TracedServeMux) HandleFunc(pattern string, fn func(http.ResponseWriter,
	*http.Request)) {
	middleware := nethttp.Middleware(
		tm.tracer,
		http.HandlerFunc(fn),
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + " " + r.URL.Path
		}),
	)
	tm.mux.Handle(pattern, middleware)
}

// ServeHTTP implements mux.Router#ServeHTTP
func (tm *TracedServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm.mux.ServeHTTP(w, r)
}
