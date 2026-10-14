// SPDX-FileCopyrightText: 2026 SunghoLee98
//
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// The game sends PKR-Client-Version with every request, so a preflight from
// another origin must allow it or the browser blocks the request.
func TestProdHandlerAllowsClientVersionHeader(t *testing.T) {
	handler := prodHandler(http.NewServeMux(), parseAllowedOrigins("https://game.example"))

	req := httptest.NewRequest(http.MethodOptions, "/account/register", nil)
	req.Header.Set("Origin", "https://game.example")
	req.Header.Set("Access-Control-Request-Headers", "authorization,content-type,pkr-client-version")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	allowed := strings.ToLower(rec.Header().Get("Access-Control-Allow-Headers"))
	for _, h := range []string{"authorization", "content-type", "pkr-client-version"} {
		if !strings.Contains(allowed, h) {
			t.Errorf("Access-Control-Allow-Headers %q does not allow %s", allowed, h)
		}
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "https://game.example" {
		t.Errorf("Access-Control-Allow-Origin = %q", got)
	}
}

// gameurl may list several game URLs (comma-separated), e.g. the GitHub Pages
// site and the Railway one. Each allowed origin gets its own origin echoed back.
func TestProdHandlerAllowsEachConfiguredOrigin(t *testing.T) {
	origins := parseAllowedOrigins("https://sungholee98.github.io/pokerogue/, https://game.example")
	handler := prodHandler(http.NewServeMux(), origins)

	for _, origin := range []string{"https://sungholee98.github.io", "https://game.example"} {
		req := httptest.NewRequest(http.MethodOptions, "/account/login", nil)
		req.Header.Set("Origin", origin)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != origin {
			t.Errorf("origin %s: Access-Control-Allow-Origin = %q", origin, got)
		}
		if got := rec.Header().Get("Vary"); !strings.Contains(got, "Origin") {
			t.Errorf("origin %s: Vary = %q, want it to include Origin", origin, got)
		}
	}
}

func TestProdHandlerRejectsUnlistedOrigin(t *testing.T) {
	handler := prodHandler(http.NewServeMux(), parseAllowedOrigins("https://game.example"))

	req := httptest.NewRequest(http.MethodOptions, "/account/login", nil)
	req.Header.Set("Origin", "https://evil.example")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got == "https://evil.example" {
		t.Errorf("Access-Control-Allow-Origin echoed unlisted origin %q", got)
	}
}

func TestParseAllowedOrigins(t *testing.T) {
	got := parseAllowedOrigins(" https://a.example/game/ ,https://b.example:8443,, https://a.example")
	want := []string{"https://a.example", "https://b.example:8443"}
	if strings.Join(got, " ") != strings.Join(want, " ") {
		t.Errorf("parseAllowedOrigins = %q, want %q", got, want)
	}
}

// OAuth callbacks redirect players back to the game, so they need a full URL:
// the first entry of gameurl.
func TestPrimaryGameURL(t *testing.T) {
	if got := primaryGameURL(" https://sungholee98.github.io/pokerogue/ , https://game.example"); got != "https://sungholee98.github.io/pokerogue/" {
		t.Errorf("primaryGameURL = %q", got)
	}
}
