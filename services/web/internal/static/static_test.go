package static

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func testFS() fstest.MapFS {
	return fstest.MapFS{
		"index.html":           {Data: []byte("<html>app</html>")},
		"assets/app-abc123.js": {Data: []byte("console.log('hi')")},
		"manifest.webmanifest": {Data: []byte("{}")},
	}
}

func get(t *testing.T, h http.Handler, method, target string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	return rec
}

func TestServesIndexAtRoot(t *testing.T) {
	rec := get(t, Handler(testFS()), http.MethodGet, "/")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "app") {
		t.Fatalf("GET / = %d %q", rec.Code, rec.Body.String())
	}
	if cc := rec.Header().Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("Cache-Control = %q, want no-cache", cc)
	}
}

func TestAppRouteFallsBackToIndex(t *testing.T) {
	rec := get(t, Handler(testFS()), http.MethodGet, "/join/r7k2-miles")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "app") {
		t.Fatalf("GET /join/r7k2-miles = %d %q, want index.html", rec.Code, rec.Body.String())
	}
}

func TestMissingAssetIs404(t *testing.T) {
	rec := get(t, Handler(testFS()), http.MethodGet, "/assets/missing.js")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("GET missing asset = %d, want 404", rec.Code)
	}
}

func TestHashedAssetIsImmutable(t *testing.T) {
	rec := get(t, Handler(testFS()), http.MethodGet, "/assets/app-abc123.js")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET asset = %d", rec.Code)
	}
	if cc := rec.Header().Get("Cache-Control"); !strings.Contains(cc, "immutable") {
		t.Errorf("Cache-Control = %q, want immutable", cc)
	}
}

func TestOtherFileIsNotCachedForever(t *testing.T) {
	rec := get(t, Handler(testFS()), http.MethodGet, "/manifest.webmanifest")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET manifest = %d", rec.Code)
	}
	if cc := rec.Header().Get("Cache-Control"); cc != "no-cache" {
		t.Errorf("Cache-Control = %q, want no-cache", cc)
	}
}

func TestRejectsWrites(t *testing.T) {
	rec := get(t, Handler(testFS()), http.MethodPost, "/")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("POST / = %d, want 405", rec.Code)
	}
}

func TestEmbeddedPlaceholderExists(t *testing.T) {
	rec := get(t, Handler(Dist()), http.MethodGet, "/")
	if rec.Code != http.StatusOK || !strings.Contains(rec.Body.String(), "Momentum") {
		t.Fatalf("embedded index = %d %q", rec.Code, rec.Body.String())
	}
}
