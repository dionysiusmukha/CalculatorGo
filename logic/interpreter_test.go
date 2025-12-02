package logic

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestInterpret_SimpleMath(t *testing.T) {
	vars := map[string]interface{}{}
	out, err := Interpret("2+2", vars)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "4" {
		t.Fatalf("expected 4, got %q", out)
	}
}

func TestInterpret_AssignAndUseVar(t *testing.T) {
	vars := map[string]interface{}{}

	out, err := Interpret("x = 10", vars)
	if err != nil {
		t.Fatalf("assign error: %v", err)
	}
	if out != "x = 10" {
		t.Fatalf("unexpected output: %q", out)
	}

	out, err = Interpret("x", vars)
	if err != nil {
		t.Fatalf("read var error: %v", err)
	}
	if strings.TrimSpace(out) != "10" {
		t.Fatalf("expected 10, got %q", out)
	}
}

func TestInterpret_CurlCommand(t *testing.T) {
	// поднимаем фейковый HTTP-сервер
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from test"))
	}))
	defer srv.Close()

	vars := map[string]interface{}{}
	cmd := "curl " + srv.URL

	out, err := Interpret(cmd, vars)
	if err != nil {
		t.Fatalf("Interpret curl error: %v", err)
	}

	// проверяем, что last_curl сохранился
	val, ok := vars["last_curl"]
	if !ok {
		t.Fatalf("expected vars[\"last_curl\"] to be set")
	}
	body, _ := val.(string)
	if body != "hello from test" {
		t.Fatalf("unexpected body: %q", body)
	}

	if !strings.Contains(out, "HTML got with") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestInterpret_CallCommand(t *testing.T) {
	vars := map[string]interface{}{}

	// перехватываем URL, чтобы не открывать браузер
	var capturedURL string
	old := openInBrowser
	openInBrowser = func(url, preferred string) error {
		capturedURL = url
		return nil
	}
	defer func() { openInBrowser = old }()

	out, err := Interpret("позвонить rita", vars)
	if err != nil {
		t.Fatalf("Interpret call error: %v", err)
	}
	if !strings.Contains(out, "rita") {
		t.Fatalf("unexpected output: %q", out)
	}
	if capturedURL == "" {
		t.Fatalf("expected browser to be called with URL, got empty")
	}
	
}
