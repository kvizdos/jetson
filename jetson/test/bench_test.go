package jetson_test

import (
	"testing"

	"github.com/kvizdos/jetson/jetson"
)

func BenchmarkJetsonScan(b *testing.B) {
	scanner := jetson.NewJetsonScanner("../../demo/arxiv-metadata-oai-snapshot.json", "categories", "hep")

	err := scanner.ReadFile()
	if err != nil {
		b.Fatal(err)
	}
	b.Log("Scanned")
	b.ResetTimer()
	for b.Loop() {
		_, err := scanner.Scan()
		if err != nil {
			b.Fatal(err)
		}
	}
}
