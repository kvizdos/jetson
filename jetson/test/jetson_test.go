package jetson_test

import (
	"testing"

	"github.com/kvizdos/jetson/jetson"
)

func newInMemoryScanner(data string, key string, value string, workers int) *jetson.LineScanner {
	return &jetson.LineScanner{
		Data:        []byte(data),
		SearchKey:   []byte(`"` + key + `":"`),
		SearchValue: []byte(value),
		Workers:     workers,
	}
}

func TestJetsonScanner_BasicMatch(t *testing.T) {
	scanner := newInMemoryScanner(
		`{"id":1,"categories":"hep-th"}
{"id":2,"categories":"math"}
{"id":3,"categories":"hep-ex"}`,
		"categories", "hep", 3)

	result, err := scanner.Scan()
	if err != nil {
		t.Fatal(err)
	}
	if result.Matches != 2 {
		t.Errorf("Expected 2 matches, got %d", result.Matches)
	}
}

func TestJetsonScanner_NoMatch(t *testing.T) {
	scanner := newInMemoryScanner(
		`{"categories":"bio"}`, "categories", "hep", 2)

	result, err := scanner.Scan()
	if err != nil {
		t.Fatal(err)
	}
	if result.Matches != 0 {
		t.Errorf("Expected 0 matches, got %d", result.Matches)
	}
}

func TestJetsonScanner_MissingKey(t *testing.T) {
	scanner := newInMemoryScanner(
		`{"other":"value"}`, "categories", "hep", 2)

	result, err := scanner.Scan()
	if err != nil {
		t.Fatal(err)
	}
	if result.Matches != 0 {
		t.Errorf("Expected 0 matches, got %d", result.Matches)
	}
}

func TestJetsonScanner_EmptyInput(t *testing.T) {
	scanner := newInMemoryScanner("", "categories", "hep", 2)

	result, err := scanner.Scan()
	if err != nil {
		t.Fatal(err)
	}
	if result.Matches != 0 {
		t.Errorf("Expected 0 matches on empty input, got %d", result.Matches)
	}
}

func TestJetsonScanner_MalformedLine(t *testing.T) {
	scanner := newInMemoryScanner(
		`{"categories":"hep"}
not_json_at_all
{"categories":"hep-ex"}`,
		"categories", "hep", 4)

	result, err := scanner.Scan()
	if err != nil {
		t.Fatal(err)
	}
	if result.Matches != 2 {
		t.Errorf("Expected 2 matches despite malformed line, got %d", result.Matches)
	}
}

func TestJetsonScanner_MoreWorkersThanLines(t *testing.T) {
	scanner := newInMemoryScanner(
		`{"categories":"hep"}`, "categories", "hep", 8)

	result, err := scanner.Scan()
	if err != nil {
		t.Fatal(err)
	}
	if result.Matches != 1 {
		t.Errorf("Expected 1 match, got %d", result.Matches)
	}
}
