package osrun

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsVideoFile(t *testing.T) {
	cases := []struct {
		name string
		want bool
	}{
		{"movie.mp4", true},
		{"clip.MKV", true},
		{"image.png", false},
		{"text.txt", false},
	}
	for _, c := range cases {
		if got := isVideoFile(c.name); got != c.want {
			t.Errorf("isVideoFile(%q) = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestDepthExceeds(t *testing.T) {
	root := "/home/user/Videos"
	if depthExceeds(root, root, 2) {
		t.Errorf("root should not exceed depth")
	}
	if !depthExceeds(root, filepath.Join(root, "a", "b", "c"), 2) {
		t.Errorf("deep path should exceed depth")
	}
}

func TestFindMediaFile_InSafeDirs(t *testing.T) {
	tmp := t.TempDir()

	// подменяем HOME так, чтобы safeDirs() смотрел в наш temp
	t.Setenv("HOME", tmp)

	// создаём Videos с файлом
	videos := filepath.Join(tmp, "Videos")
	if err := os.MkdirAll(videos, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	fpath := filepath.Join(videos, "test_video.mp4")
	if err := os.WriteFile(fpath, []byte("dummy"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	got, err := FindMediaFile("test_video.mp4", 3)
	if err != nil {
		t.Fatalf("FindMediaFile error: %v", err)
	}
	if got != fpath {
		t.Fatalf("expected %q, got %q", fpath, got)
	}
}

func TestOpenInBrowserPrefer_BadURL(t *testing.T) {
	if err := OpenInBrowserPrefer("not-a-url", ""); err == nil {
		t.Fatalf("expected error on bad url, got nil")
	}
}
