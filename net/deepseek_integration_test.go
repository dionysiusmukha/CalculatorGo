//go:build integration

package net

import "testing"

func TestSendToDeepSeek_Integration(t *testing.T) {
	resp, err := SendToDeepSeek("Скажи число 42 и только его.")
	if err != nil {
		t.Fatalf("SendToDeepSeek error: %v", err)
	}
	if resp == "" {
		t.Fatalf("empty response from DeepSeek")
	}
}

func TestExtractFreeForm_Integration(t *testing.T) {
	ex, err := ExtractFreeForm("открой сайт https://example.com в браузере")
	if err != nil {
		t.Fatalf("ExtractFreeForm error: %v", err)
	}
	if ex.Action == "" {
		t.Fatalf("empty action in Extracted")
	}
}
