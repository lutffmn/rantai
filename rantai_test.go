package rantai

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockMiddleware(name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Middleware", name)
			next.ServeHTTP(w, r)
		})
	}
}

func mockHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func TestNew(t *testing.T) {
	mw1 := mockMiddleware("mw1")
	mw2 := mockMiddleware("mw2")
	r := New(mw1, mw2)
	if len(r.Middlewares) != 2 {
		t.Errorf("expected 2 middlewares, got %d", len(r.Middlewares))
	}
}

func TestChain(t *testing.T) {
	mw1 := mockMiddleware("mw1")
	mw2 := mockMiddleware("mw2")
	r := New(mw1, mw2)
	handler := r.Chain(http.HandlerFunc(mockHandler))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Body.String() != "OK" {
		t.Errorf("expected body 'OK', got %q", w.Body.String())
	}
	headers := w.Header().Get("X-Middleware")
	if headers != "mw2" {
		t.Errorf("expected header 'mw2', got %q", headers)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil handler")
		}
	}()
	r.Chain(nil)
}

func TestChainFunc(t *testing.T) {
	mw1 := mockMiddleware("mw1")
	mw2 := mockMiddleware("mw2")
	r := New(mw1, mw2)
	handler := r.ChainF(mockHandler)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Body.String() != "OK" {
		t.Errorf("expected body 'OK', got %q", w.Body.String())
	}
	headers := w.Header().Get("X-Middleware")
	if headers != "mw2" {
		t.Errorf("expected header 'mw2', got %q", headers)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil handler")
		}
	}()
	r.Chain(nil)
}

func TestExtend(t *testing.T) {
	mw1 := mockMiddleware("mw1")
	r := New(mw1)
	mw2 := mockMiddleware("mw2")
	r2 := r.Extend(mw2)

	if len(r.Middlewares) != 1 {
		t.Errorf("original Rantai should have 1 middleware, got %d", len(r.Middlewares))
	}
	if len(r2.Middlewares) != 2 {
		t.Errorf("extended Rantai should have 2 middlewares, got %d", len(r2.Middlewares))
	}
}

func TestExclude(t *testing.T) {
	mw1 := mockMiddleware("mw1")
	mw2 := mockMiddleware("mw2")
	r := New(mw1, mw2)
	r2 := r.Exclude(mw1)

	if len(r2.Middlewares) != 1 {
		t.Errorf("expected 1 middleware after exclusion, got %d", len(r2.Middlewares))
	}

	handler := r2.Chain(http.HandlerFunc(mockHandler))
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	headers := w.Header().Get("X-Middleware")
	if headers != "mw2" {
		t.Errorf("expected header 'mw2', got %q", headers)
	}
}

