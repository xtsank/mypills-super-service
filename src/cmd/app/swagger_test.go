package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSwaggerRoutes(t *testing.T) {
	root, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatal(err)
	}
	previous, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previous) })

	app := &App{router: gin.New()}
	app.initSwagger()

	for _, path := range []string{"/swagger/", "/swagger/swagger-initializer.js",
		"/swagger/swagger-ui-bundle.js", "/swagger-v1/index.html",
		"/swagger-v1/swagger-initializer.js", "/swagger-v1/swagger-ui-bundle.js",
		"/swagger/swagger-ui-bundle.js", "/api/spec/v2.yaml"} {
		recorder := httptest.NewRecorder()
		app.router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusOK {
			t.Errorf("%s: status %d, body %q", path, recorder.Code, recorder.Body.String())
		}
		if path == "/api/spec/v2.yaml" && !strings.Contains(recorder.Body.String(), "openapi: 3.1.0") {
			t.Errorf("%s: missing OpenAPI document", path)
		}
		if path == "/swagger/swagger-initializer.js" && !strings.Contains(recorder.Body.String(), "/api/spec/v2.yaml") {
			t.Errorf("%s: missing v2 specification URL", path)
		}
	}

	for _, path := range []string{"/swagger", "/swagger-v1/"} {
		recorder := httptest.NewRecorder()
		app.router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))
		if recorder.Code != http.StatusFound {
			t.Errorf("%s: expected redirect, got %d", path, recorder.Code)
		}
	}
}
