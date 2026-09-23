package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootUsesRealClientIP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "203.0.113.10, 198.51.100.5")
	rec := httptest.NewRecorder()

	New().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if body := rec.Body.String(); !strings.Contains(body, "203.0.113.10") {
		t.Errorf("response body does not contain forwarded client IP: %q", body)
	}
	for _, want := range []string{
		`src="/assets/tailwindcss-browser-v4.3.3.js"`,
		`type="module" src="/assets/datastar-v1.0.4.js"`,
		`data-on:click="$showHeaders = !$showHeaders"`,
		`data-show="$showHeaders"`,
		`data-on:click="@get('/server-time')"`,
		`id="server-time"`,
	} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("root page does not contain %q", want)
		}
	}
}

func TestTailwindAsset(t *testing.T) {
	rec := httptest.NewRecorder()
	New().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/assets/tailwindcss-browser-v4.3.3.js", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "javascript") {
		t.Errorf("Content-Type = %q, want JavaScript", ct)
	}
	if !strings.Contains(rec.Body.String(), "Minified by jsDelivr") {
		t.Error("Tailwind asset body is missing")
	}
}

func TestDatastarAsset(t *testing.T) {
	rec := httptest.NewRecorder()
	New().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/assets/datastar-v1.0.4.js", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.Contains(ct, "javascript") {
		t.Errorf("Content-Type = %q, want JavaScript", ct)
	}
	if !strings.Contains(rec.Body.String(), "Datastar v1.0.4") {
		t.Error("Datastar asset body is missing")
	}
}

func TestServerTime(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/server-time", nil)
	req.Header.Set("Accept", "text/event-stream")
	New().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Errorf("Content-Type = %q, want text/event-stream", ct)
	}
	for _, want := range []string{"event: datastar-patch-elements", `id="server-time"`, "Server time: "} {
		if !strings.Contains(rec.Body.String(), want) {
			t.Errorf("SSE response does not contain %q: %q", want, rec.Body.String())
		}
	}
	if strings.Contains(rec.Body.String(), "Not requested yet") {
		t.Error("SSE response still contains the initial placeholder")
	}
}

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

			New().ServeHTTP(rec, req)

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
