package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func LoadVars() map[string]interface{} {
	var vars map[string]interface{}
	b, err := os.ReadFile("vars.json")
	if err == nil {
		if err := json.Unmarshal(b, &vars); err != nil {
			fmt.Println("Warning: could not parse vars.json:", err)
		}
	}
	if vars == nil {
		vars = make(map[string]interface{})
	}
	return vars
}

func SaveVars(vars map[string]interface{}) {
	data, _ := json.MarshalIndent(vars, "", "  ")
	if err := os.WriteFile("vars.json", data, 0644); err != nil {
		fmt.Println("Error writing vars.json:", err)
	}
}

func LoadHistory() []string {
	var history []string
	b, err := os.ReadFile("history.txt")
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(b)), "\n")
		if len(lines) > 10 {
			lines = lines[len(lines)-10:]
		}
		history = append(history, lines...)
	}
	return history
}

func SaveHistory(history []string) {
	_ = os.WriteFile("history.txt", []byte(strings.Join(history, "\n")), 0644)
}
