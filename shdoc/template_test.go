package shdoc

import (
	"os"
	"testing"
)

func renderExample(t *testing.T, example string, trimTabs int) string {
	t.Helper()

	doc := Document{
		Sections: []Section{{Functions: []FuncDoc{{Examples: []string{example}}}}},
	}
	out, err := RenderWithTemplateOptions(&doc, `{{formatExample (index (index (index .Sections 0).Functions 0).Examples 0)}}`, RenderOptions{
		ExampleTrimTabs: trimTabs,
	})
	if err != nil {
		t.Fatalf("RenderWithTemplateOptions() error: %v", err)
	}
	return out
}

func TestCustomTemplate(t *testing.T) {
	input := "# @name mylib\n"
	doc, _ := ParseDocument(input)
	out, err := RenderWithTemplate(&doc, `{{.FileTitle}}`)
	if err != nil {
		t.Fatalf("renderWithTemplate error: %v", err)
	}
	if out != "mylib" {
		t.Errorf("got %q, want %q", out, "mylib")
	}
}

func TestPrintTemplateRoundtrip(t *testing.T) {
	input, err := os.ReadFile("../examples/showcase.sh")
	if err != nil {
		t.Fatalf("open showcase.sh: %v", err)
	}
	doc, _ := ParseDocument(string(input))
	out1, err := RenderWithTemplate(&doc, DefaultMarkdownTemplate)
	if err != nil {
		t.Fatalf("renderWithTemplate error: %v", err)
	}
	// Render again with the same template to verify determinism.
	out2, err := RenderWithTemplate(&doc, DefaultMarkdownTemplate)
	if err != nil {
		t.Fatalf("renderWithTemplate (second call) error: %v", err)
	}
	if out1 != out2 {
		t.Errorf("two renders of the same document differ")
	}
}

func TestRenderWithTemplateOptionsFormatsExamplesWithTabTrim(t *testing.T) {
	out := renderExample(t, "\t\tpodman run\n \tleave-this\n\tone-tab\n\t\tsecond-line", 2)
	want := "podman run\n \tleave-this\none-tab\nsecond-line"
	if out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

func TestRenderWithTemplateOptionsSkipsTabTrimWhenFirstLineDoesNotMeetThreshold(t *testing.T) {
	out := renderExample(t, "\tpodman run\n\t\tsecond-line", 2)
	if out != "\tpodman run\n\t\tsecond-line" {
		t.Errorf("got %q, want %q", out, "\tpodman run\n\t\tsecond-line")
	}
}

func TestRenderWithTemplateOptionsDisablesExampleTabTrimAtZero(t *testing.T) {
	out := renderExample(t, "\t\tpodman run", 0)
	if out != "\t\tpodman run" {
		t.Errorf("got %q, want %q", out, "\t\tpodman run")
	}
}
