package static

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/aisk/vox"
)

func TestMiddleware(t *testing.T) {
	root := t.TempDir()
	err := os.WriteFile(filepath.Join(root, "hello.txt"), []byte("Hello Static!"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	app := vox.New()
	app.SetConfig("logging:disable", "true")
	app.Use(Middleware("/public", root))

	r := httptest.NewRequest("GET", "http://test.com/public/hello.txt", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expect StatusCode 200, got %d", w.Result().StatusCode)
	}
	body, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "Hello Static!" {
		t.Fatalf("expect body %q, got %q", "Hello Static!", string(body))
	}

	r = httptest.NewRequest("HEAD", "http://test.com/public/hello.txt", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expect StatusCode 200, got %d", w.Result().StatusCode)
	}
	body, err = io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) != 0 {
		t.Fatalf("expect empty body for HEAD, got %q", string(body))
	}

	r = httptest.NewRequest("GET", "http://test.com/public/not-found.txt", nil)
	w = httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("expect StatusCode 404, got %d", w.Result().StatusCode)
	}
}

func TestMiddlewareNormalizePrefix(t *testing.T) {
	root := t.TempDir()
	err := os.WriteFile(filepath.Join(root, "a.txt"), []byte("A"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	app := vox.New()
	app.SetConfig("logging:disable", "true")
	app.Use(Middleware("assets/", root))

	r := httptest.NewRequest("GET", "http://test.com/assets/a.txt", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expect StatusCode 200, got %d", w.Result().StatusCode)
	}
}

func TestMiddlewareWithRootPrefix(t *testing.T) {
	root := t.TempDir()
	err := os.WriteFile(filepath.Join(root, "root.txt"), []byte("Root"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	app := vox.New()
	app.SetConfig("logging:disable", "true")
	app.Use(Middleware("", root))

	r := httptest.NewRequest("GET", "http://test.com/root.txt", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expect StatusCode 200, got %d", w.Result().StatusCode)
	}
	body, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "Root" {
		t.Fatalf("expect body %q, got %q", "Root", string(body))
	}
}

func TestMiddlewareWillNotOverrideExistingResponse(t *testing.T) {
	root := t.TempDir()
	err := os.WriteFile(filepath.Join(root, "hello.txt"), []byte("Hello Static!"), 0o644)
	if err != nil {
		t.Fatal(err)
	}

	app := vox.New()
	app.SetConfig("logging:disable", "true")
	app.Get("/public/hello.txt", func(ctx *vox.Context, req *vox.Request[vox.NoBody], res *vox.Response[any]) {
		res.Body = "from route"
	})
	app.Use(Middleware("/public", root))

	r := httptest.NewRequest("GET", "http://test.com/public/hello.txt", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expect StatusCode 200, got %d", w.Result().StatusCode)
	}
	body, err := io.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "from route" {
		t.Fatalf("expect body %q, got %q", "from route", string(body))
	}
}
