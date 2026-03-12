package millerutils

import (
	"strings"
	"testing"
)

// sample N-Quads with blank nodes that have a schema:contentUrl property
const nqWithContentURL = `_:b0 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/DataDownload> .
_:b0 <http://schema.org/contentUrl> "http://example.com/data.csv" .
_:b0 <http://schema.org/encodingFormat> "text/csv" .
`

// sample N-Quads with blank nodes that have schema:name + rdf:type
const nqWithName = `_:b0 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/Person> .
_:b0 <http://schema.org/name> "Jane Doe" .
_:b0 <http://schema.org/email> "jane@example.com" .
`

// sample N-Quads with blank nodes that have schema:title + rdf:type
const nqWithTitle = `_:b0 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/CreativeWork> .
_:b0 <http://schema.org/title> "My Great Work" .
`

// sample N-Quads where blank node has no identifying properties (only generic predicates)
const nqGenericProps = `_:b0 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/Thing> .
_:b0 <http://schema.org/description> "something" .
`

// sample N-Quads with multiple blank nodes
const nqMultiple = `<http://example.org/dataset1> <http://schema.org/distribution> _:b0 .
<http://example.org/dataset1> <http://schema.org/creator> _:b1 .
_:b0 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/DataDownload> .
_:b0 <http://schema.org/contentUrl> "http://example.com/data.csv" .
_:b1 <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://schema.org/Person> .
_:b1 <http://schema.org/name> "John Smith" .
`

// blank node with no properties at all (only referenced as object)
const nqNoProps = `<http://example.org/thing> <http://schema.org/about> _:b0 .
`

func TestDeterministicBNodes_Idempotent(t *testing.T) {
	sourceURL := "http://example.org/page1"

	result1 := DeterministicBNodes(nqWithContentURL, sourceURL)
	result2 := DeterministicBNodes(nqWithContentURL, sourceURL)

	if result1 != result2 {
		t.Errorf("DeterministicBNodes is not idempotent.\nFirst:  %s\nSecond: %s", result1, result2)
	}
}

func TestDeterministicBNodes_NoOriginalBNodes(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqWithContentURL, sourceURL)

	if strings.Contains(result, "_:b0") {
		t.Error("original blank node _:b0 was not replaced")
	}
}

func TestDeterministicBNodes_DifferentSourceURLs(t *testing.T) {
	result1 := DeterministicBNodes(nqWithContentURL, "http://example.org/page1")
	result2 := DeterministicBNodes(nqWithContentURL, "http://example.org/page2")

	// Extract the generated bnode IDs
	id1 := extractBNodeID(result1)
	id2 := extractBNodeID(result2)

	if id1 == id2 {
		t.Errorf("different source URLs should produce different bnode IDs, got same: %s", id1)
	}
}

func TestDeterministicBNodes_ContentURLPriority(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqWithContentURL, sourceURL)

	// Should be deterministic and not contain original _:b0
	if strings.Contains(result, "_:b0") {
		t.Error("blank node was not replaced")
	}

	// Verify it's a hash-based ID (starts with _:b followed by hex)
	id := extractBNodeID(result)
	if !strings.HasPrefix(id, "_:b") {
		t.Errorf("expected _:b prefix, got: %s", id)
	}
	if len(id) < 10 {
		t.Errorf("expected longer deterministic ID, got: %s", id)
	}
}

func TestDeterministicBNodes_NameFallback(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqWithName, sourceURL)

	if strings.Contains(result, "_:b0") {
		t.Error("blank node was not replaced using name fallback")
	}
}

func TestDeterministicBNodes_TitleFallback(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqWithTitle, sourceURL)

	if strings.Contains(result, "_:b0") {
		t.Error("blank node was not replaced using title fallback")
	}
}

func TestDeterministicBNodes_GenericPropsFallback(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqGenericProps, sourceURL)

	if strings.Contains(result, "_:b0") {
		t.Error("blank node was not replaced using generic properties fallback")
	}

	// Should still be deterministic
	result2 := DeterministicBNodes(nqGenericProps, sourceURL)
	if result != result2 {
		t.Error("generic property fallback should be deterministic")
	}
}

func TestDeterministicBNodes_MultipleBNodes(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqMultiple, sourceURL)

	if strings.Contains(result, "_:b0") || strings.Contains(result, "_:b1") {
		t.Error("not all blank nodes were replaced")
	}

	// The two different blank nodes should get different IDs
	lines := strings.Split(strings.TrimSpace(result), "\n")
	bnodes := make(map[string]bool)
	for _, line := range lines {
		parts := strings.Split(line, " ")
		for _, p := range parts {
			if strings.HasPrefix(p, "_:b") {
				bnodes[p] = true
			}
		}
	}
	if len(bnodes) < 2 {
		t.Errorf("expected at least 2 distinct bnode IDs, got %d: %v", len(bnodes), bnodes)
	}
}

func TestDeterministicBNodes_NoPropsNode(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqNoProps, sourceURL)

	// Should still replace the blank node (with random fallback)
	if strings.Contains(result, "_:b0") {
		t.Error("blank node with no properties was not replaced")
	}
}

func TestDeterministicBNodes_EmptyInput(t *testing.T) {
	result := DeterministicBNodes("", "http://example.org/page1")
	if result != "" {
		t.Errorf("expected empty output for empty input, got: %s", result)
	}
}

func TestUnwrapURI(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"<http://schema.org/name>", "http://schema.org/name"},
		{"http://schema.org/name", "http://schema.org/name"},
		{"<>", ""},
		{"plain", "plain"},
	}
	for _, tc := range tests {
		got := unwrapURI(tc.input)
		if got != tc.expected {
			t.Errorf("unwrapURI(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestGlobalUniqueBNodes_StillWorks(t *testing.T) {
	// Ensure the old function still works as a fallback
	result := GlobalUniqueBNodes(nqWithContentURL)
	if strings.Contains(result, "_:b0") {
		t.Error("GlobalUniqueBNodes did not replace blank nodes")
	}
}

// extractBNodeID extracts the first blank node ID from an N-Quads string
func extractBNodeID(nq string) string {
	for _, line := range strings.Split(nq, "\n") {
		parts := strings.Split(line, " ")
		for _, p := range parts {
			if strings.HasPrefix(p, "_:b") && len(p) > 4 {
				return p
			}
		}
	}
	return ""
}
