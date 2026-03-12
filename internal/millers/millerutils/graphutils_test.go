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

func TestDeterministicBNodes_ProducesURIs(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqWithContentURL, sourceURL)

	if !strings.Contains(result, GleanerBaseURI) {
		t.Errorf("expected skolemized URI with base %s, got:\n%s", GleanerBaseURI, result)
	}
}

func TestDeterministicBNodes_ContentURLIncludesType(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqWithContentURL, sourceURL)

	if !strings.Contains(result, "/DataDownload/") {
		t.Errorf("expected type 'DataDownload' in URI, got:\n%s", result)
	}
}

func TestDeterministicBNodes_NameIncludesType(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqWithName, sourceURL)

	if !strings.Contains(result, "/Person/") {
		t.Errorf("expected type 'Person' in URI, got:\n%s", result)
	}
}

func TestDeterministicBNodes_TitleIncludesType(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqWithTitle, sourceURL)

	if !strings.Contains(result, "/CreativeWork/") {
		t.Errorf("expected type 'CreativeWork' in URI, got:\n%s", result)
	}
}

func TestDeterministicBNodes_DifferentSourceURLs(t *testing.T) {
	result1 := DeterministicBNodes(nqWithContentURL, "http://example.org/page1")
	result2 := DeterministicBNodes(nqWithContentURL, "http://example.org/page2")

	// Extract the generated URIs
	uri1 := extractGenID(result1)
	uri2 := extractGenID(result2)

	if uri1 == uri2 {
		t.Errorf("different source URLs should produce different URIs, got same: %s", uri1)
	}
}

func TestDeterministicBNodes_GenericPropsFallback(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqGenericProps, sourceURL)

	if strings.Contains(result, "_:b0") {
		t.Error("blank node was not replaced using generic properties fallback")
	}

	if !strings.Contains(result, "/Thing/") {
		t.Errorf("expected type 'Thing' in fallback URI, got:\n%s", result)
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

	// Should contain both type segments
	if !strings.Contains(result, "/DataDownload/") {
		t.Error("missing DataDownload type in URI")
	}
	if !strings.Contains(result, "/Person/") {
		t.Error("missing Person type in URI")
	}
}

func TestDeterministicBNodes_NoPropsNode(t *testing.T) {
	sourceURL := "http://example.org/page1"
	result := DeterministicBNodes(nqNoProps, sourceURL)

	// Should still replace the blank node
	if strings.Contains(result, "_:b0") {
		t.Error("blank node with no properties was not replaced")
	}

	// Node with no type gets "Unknown"
	if !strings.Contains(result, "/Unknown/") {
		t.Errorf("expected 'Unknown' type for property-less node, got:\n%s", result)
	}
}

func TestDeterministicBNodes_EmptyInput(t *testing.T) {
	result := DeterministicBNodes("", "http://example.org/page1")
	if result != "" {
		t.Errorf("expected empty output for empty input, got: %s", result)
	}
}

func TestExtractTypeName(t *testing.T) {
	tests := []struct {
		input, expected string
	}{
		{"<http://schema.org/DataDownload>", "DataDownload"},
		{"<https://schema.org/Person>", "Person"},
		{"<http://www.w3.org/1999/02/22-rdf-syntax-ns#Class>", "Class"},
		{"", "Unknown"},
		{"SomePlainValue", "SomePlainValue"},
	}
	for _, tc := range tests {
		got := extractTypeName(tc.input)
		if got != tc.expected {
			t.Errorf("extractTypeName(%q) = %q, want %q", tc.input, got, tc.expected)
		}
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

// extractGenID extracts the first gleaner genid URI from an N-Quads string
func extractGenID(nq string) string {
	start := strings.Index(nq, "<"+GleanerBaseURI)
	if start < 0 {
		return ""
	}
	end := strings.Index(nq[start:], ">")
	if end < 0 {
		return ""
	}
	return nq[start : start+end+1]
}
