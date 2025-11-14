package osrun

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var videoExts = []string{".mp4", ".mkv", ".avi", ".mov", ".webm"}
var browserCandidates = []string{
	"yandex-browser-stable", "yandex-browser",
	"firefox",
	"google-chrome-stable", "google-chrome",
	"chromium",
	"brave",
	"xdg-open",
}
var playerCandidates  = []string{"mpv", "xdg-open"}

func normalizeBrowserName(app string) []string {
	a := strings.ToLower(strings.TrimSpace(app))
	switch a {
	case "yandex", "yabrowser", "yandex-browser", "yandex browser", "яндекс", "яндекс-браузер", "yandex-browser-stable":
		return []string{"yandex-browser-stable", "yandex-browser"}
	case "chrome", "google-chrome":
		return []string{"google-chrome-stable", "google-chrome"}
	case "firefox", "ff":
		return []string{"firefox"}
	case "chromium":
		return []string{"chromium"}
	case "brave":
		return []string{"brave"}
	default:
		if a != "" {
			return []string{a}
		}
		return nil
	}
}


func safeDirs() ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return []string{
		filepath.Join(home, "Videos"),
		filepath.Join(home, "Downloads"),
		filepath.Join(home, "Music"),
	}, nil
}


func isSafePath(p string) (bool, error) {
	p = filepath.Clean(p)
	resolved, err := filepath.EvalSymlinks(p)
	if err != nil {
		return false, err
	}
	allowed, err := safeDirs()
	if err != nil {
		return false, err
	}
	for _, base := range allowed {
		baseRes, err := filepath.EvalSymlinks(base)
		if err == nil {
			if strings.HasPrefix(resolved, baseRes+string(filepath.Separator)) || resolved == baseRes {
				return true, nil
			}
		}
	}
	return false, nil
}

func isVideoFile(name string) bool {
	low := strings.ToLower(name)
	for _, ext := range videoExts {
		if strings.HasSuffix(low, ext) {
			return true
		}
	}
	return false
}


func FindMediaFile(filename string, maxDepth int) (string, error) {
	if filename == "" {
		return "", errors.New("empty filename")
	}
	if !isVideoFile(filename) {
		cands := []string{filename}
		for _, ext := range videoExts {
			if !strings.HasSuffix(strings.ToLower(filename), ext) {
				cands = append(cands, filename+ext)
			}
		}
		filename = filename
		_ = cands
	}

	allowed, err := safeDirs()
	if err != nil {
		return "", err
	}

	filenameLow := strings.ToLower(filename)

	for _, root := range allowed {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, e error) error {
			if e != nil {
				return nil
			}
			if d.IsDir() {
				if depthExceeds(root, path, maxDepth) {
					return filepath.SkipDir
				}
				return nil
			}
			base := strings.ToLower(d.Name())
			if (base == filenameLow || strings.HasPrefix(base, filenameLow)) && isVideoFile(base) {
				ok, err := isSafePath(path)
				if err == nil && ok {
					return errors.New("FOUND::" + path)
				}
			}
			return nil
		})
		if err != nil && strings.HasPrefix(err.Error(), "FOUND::") {
			return strings.TrimPrefix(err.Error(), "FOUND::"), nil
		}
	}
	return "", fmt.Errorf("файл %q не найден в безопасных директориях", filename)
}

func depthExceeds(root, path string, maxDepth int) bool {
	if maxDepth <= 0 {
		return false
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	if rel == "." {
		return false
	}
	depth := 0
	for _, p := range strings.Split(rel, string(filepath.Separator)) {
		if p != "" {
			depth++
		}
	}
	return depth > maxDepth
}




// Открыть URL, с учётом желаемого браузера (если указан)
func OpenInBrowserPrefer(url, preferred string) error {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return fmt.Errorf("некорректный URL: %s", url)
	}

	// 1) Если пользователь просил конкретный браузер
	if prefs := normalizeBrowserName(preferred); len(prefs) > 0 {
		for _, bin := range prefs {
			if _, err := exec.LookPath(bin); err == nil {
				return exec.Command(bin, url).Start()
			}
		}
	}

	// 2) Иначе — по общему whitelist
	for _, bin := range browserCandidates {
		if _, err := exec.LookPath(bin); err == nil {
			return exec.Command(bin, url).Start()
		}
	}
	return fmt.Errorf("не найден ни один браузер из whitelist")
}


func OpenInBrowser(url string) error {
	return OpenInBrowserPrefer(url, "")
}


func OpenWithPlayer(file string) error {
	ok, err := isSafePath(file)
	if err != nil || !ok {
		return fmt.Errorf("путь небезопасен или недоступен: %v", err)
	}
	for _, bin := range playerCandidates {
		if _, err := exec.LookPath(bin); err == nil {
			cmd := exec.Command(bin, file)
			return cmd.Start()
		}
	}
	return fmt.Errorf("не найден ни один видеоплеер из whitelist")
}
