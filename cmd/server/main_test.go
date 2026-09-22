package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	tests := []struct {
		name         string
		accept       string
		acceptValues []string
		wantStatus   int
		wantCT       string
		wantBody     string
	}{
		{
			name:       "no accept header defaults to plain text",
			wantStatus: http.StatusOK,
			wantCT:     "text/plain; charset=utf-8",
			wantBody:   "OK",
		},
		{
			name:       "browser accept defaults to plain text",
			accept:     "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
			wantStatus: http.StatusOK,
			wantCT:     "text/plain; charset=utf-8",
			wantBody:   "OK",
		},
		{
			name:       "json accept returns json",
			accept:     "application/json",
			wantStatus: http.StatusOK,
			wantCT:     "application/json",
			wantBody:   `{"status":"ok"}`,
		},
		{
			name:       "json within accept list returns json",
			accept:     "text/html, application/json",
			wantStatus: http.StatusOK,
			wantCT:     "application/json",
			wantBody:   `{"status":"ok"}`,
		},
		{
			name:       "json with parameters returns json",
			accept:     "application/json; charset=utf-8; q=0.9",
			wantStatus: http.StatusOK,
			wantCT:     "application/json",
			wantBody:   `{"status":"ok"}`,
		},
		{
			name:       "json with q=0 is not accepted",
			accept:     "application/json;q=0",
			wantStatus: http.StatusOK,
			wantCT:     "text/plain; charset=utf-8",
			wantBody:   "OK",
		},
		{
			name:       "json media type is case insensitive",
			accept:     "Application/JSON",
			wantStatus: http.StatusOK,
			wantCT:     "application/json",
			wantBody:   `{"status":"ok"}`,
		},
		{
			name:         "json in second accept header field returns json",
			acceptValues: []string{"text/html", "application/json"},
			wantStatus:   http.StatusOK,
			wantCT:       "application/json",
			wantBody:     `{"status":"ok"}`,
		},
		{
			name:         "json with q=0 in one field is ignored across fields",
			acceptValues: []string{"application/json;q=0", "text/html"},
			wantStatus:   http.StatusOK,
			wantCT:       "text/plain; charset=utf-8",
			wantBody:     "OK",
		},
		{
			name:         "json with q=0 in first field still matches later field",
			acceptValues: []string{"application/json;q=0", "application/json"},
			wantStatus:   http.StatusOK,
			wantCT:       "application/json",
			wantBody:     `{"status":"ok"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			if tt.accept != "" {
				req.Header.Set("Accept", tt.accept)
			}
			for _, value := range tt.acceptValues {
				req.Header.Add("Accept", value)
			}
			rec := httptest.NewRecorder()

			newRouter().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if ct := rec.Header().Get("Content-Type"); ct != tt.wantCT {
				t.Errorf("Content-Type = %q, want %q", ct, tt.wantCT)
			}
			if vary := rec.Header().Get("Vary"); vary != "Accept" {
				t.Errorf("Vary = %q, want %q", vary, "Accept")
			}
			if body := rec.Body.String(); body != tt.wantBody {
				t.Errorf("body = %q, want %q", body, tt.wantBody)
			}
		})
	}
}
