package linter

import (
	"testing"
)

func TestReplaceTextNodes_AnchorHrefUnchangedAndTextReplaced(t *testing.T) {
	input := []byte(`<p><a href="http://example.com/?q=Sneaked">Link</a> Sneaked “quote” -- end</p>`)
	out, err := ReplaceTextNodesInXHTML(input)
	if err != nil {
		t.Fatalf("ReplaceTextNodesInXHTML error: %v", err)
	}
	output := string(out)

	// href attribute should be preserved exactly
	if !contains(output, `href="http://example.com/?q=Sneaked"`) {
		t.Fatalf("href attribute changed or missing: %s", output)
	}

	// text node "Sneaked" should be replaced with "Snuck"
	if !contains(output, "Snuck") {
		t.Fatalf("text replacement not applied: %s", output)
	}

	// smart quotes should be replaced with straight quotes
	if !contains(output, `"quote"`) {
		t.Fatalf("smart quotes not replaced: %s", output)
	}

	// double dash should be replaced with an em dash
	if !contains(output, "— end") {
		t.Fatalf("double dash not replaced with em dash: %s", output)
	}
}

func TestReplaceTextNodes_SkipScriptAndStyle(t *testing.T) {
	input := []byte(`<script>var s = "Sneaked"; // “weird”</script><style>.cls{content:"Sneaked";}</style>`)
	out, err := ReplaceTextNodesInXHTML(input)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	output := string(out)

	// script and style contents should be preserved (no replacements)
	if !contains(output, `var s = "Sneaked"; // “weird”`) {
		t.Fatalf("script content changed: %s", output)
	}
	if !contains(output, `.cls{content:"Sneaked";}`) {
		t.Fatalf("style content changed: %s", output)
	}
}

// contains is a small helper to avoid importing strings in test functions repeatedly
func contains(s, substr string) bool {
	return stringsContains(s, substr)
}

func stringsContains(s, substr string) bool { return len(s) >= len(substr) && (func() bool { return (func() bool { return stringsIndex(s, substr) >= 0 })() })() }

// minimal implementations to avoid importing strings in the test file body
func stringsIndex(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
