package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMaxBodySizeMiddleware_UnderLimit(t *testing.T) {
	s := &Server{}
	body := `{"url":"https://github.com/user/repo"}`
	handler := s.maxBodySizeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("unexpected error reading body: %v", err)
		}
		if string(data) != body {
			t.Errorf("expected body %q, got %q", body, string(data))
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/api/check", strings.NewReader(body))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 for body under limit, got %d", rr.Code)
	}
}

func TestMaxBodySizeMiddleware_OverLimit(t *testing.T) {
	s := &Server{}
	// 1 MB + 1 byte exceeds the limit
	oversized := strings.Repeat("a", (1<<20)+1)
	handler := s.maxBodySizeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := io.ReadAll(r.Body)
		if err == nil {
			t.Fatal("expected error reading oversized body, got nil")
		}
		w.WriteHeader(http.StatusRequestEntityTooLarge)
	}))

	req := httptest.NewRequest("POST", "/api/check", strings.NewReader(oversized))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected 413 for oversized body, got %d", rr.Code)
	}
}

func TestMaxBodySizeMiddleware_GETUnaffected(t *testing.T) {
	s := &Server{}
	handler := s.maxBodySizeMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/jobs", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 for GET request, got %d", rr.Code)
	}
}
