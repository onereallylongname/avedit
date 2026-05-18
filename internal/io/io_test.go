package io

import (
	"os"
	"path/filepath"
	"testing"
)

func samplesDir() string {
	return filepath.Join("..", "..", "..", "samples")
}

func TestLoadAvroFromBytes_SimpleRecord(t *testing.T) {
	input := []byte(`{
		"type": "record",
		"name": "User",
		"fields": [
			{"name": "id", "type": "int"},
			{"name": "email", "type": "string"}
		]
	}`)

	proj, raw, err := LoadAvroFromBytes(input)
	if err != nil {
		t.Fatalf("LoadAvroFromBytes failed: %v", err)
	}
	if proj == nil {
		t.Fatal("expected non-nil projection")
	}
	if raw == nil {
		t.Fatal("expected non-nil raw map")
	}
	if raw["name"] != "User" {
		t.Errorf("expected raw[name]='User', got %v", raw["name"])
	}
}

func TestLoadAvroFromBytes_InvalidJSON(t *testing.T) {
	_, _, err := LoadAvroFromBytes([]byte(`{not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestLoadAvroFromBytes_MissingType(t *testing.T) {
	_, _, err := LoadAvroFromBytes([]byte(`{"name":"Foo"}`))
	if err == nil {
		t.Fatal("expected error for missing type field")
	}
}

func TestLoadAvroFromFile_Sample(t *testing.T) {
	path := filepath.Join(samplesDir(), "user-event.avsc")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("sample file not found: %s", path)
	}

	proj, raw, err := LoadAvroFromFile(path)
	if err != nil {
		t.Fatalf("LoadAvroFromFile failed: %v", err)
	}
	if proj == nil {
		t.Fatal("expected non-nil projection")
	}
	if raw["name"] == nil {
		t.Error("expected raw map to have 'name' key")
	}
}

func TestRoundtrip_LoadAndExport(t *testing.T) {
	input := []byte(`{
		"type": "record",
		"name": "Event",
		"namespace": "com.example",
		"fields": [
			{"name": "ts", "type": "long"},
			{"name": "payload", "type": "bytes"}
		]
	}`)

	proj, _, err := LoadAvroFromBytes(input)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "output.avsc")

	if err := ExportAvroToFile(proj, outPath); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Re-load exported file
	proj2, raw2, err := LoadAvroFromFile(outPath)
	if err != nil {
		t.Fatalf("re-load failed: %v", err)
	}
	if proj2 == nil {
		t.Fatal("expected non-nil re-loaded projection")
	}
	if raw2["name"] != "Event" {
		t.Errorf("expected name='Event', got %v", raw2["name"])
	}
	if raw2["namespace"] != "com.example" {
		t.Errorf("expected namespace='com.example', got %v", raw2["namespace"])
	}

	// Check fields preserved
	fields, ok := raw2["fields"].([]any)
	if !ok {
		t.Fatal("expected fields array")
	}
	if len(fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(fields))
	}
}

func TestSerializeAvro_Pretty(t *testing.T) {
	input := []byte(`{"type":"record","name":"T","fields":[]}`)
	proj, _, err := LoadAvroFromBytes(input)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	pretty, err := SerializeAvro(proj, true)
	if err != nil {
		t.Fatalf("serialize pretty failed: %v", err)
	}
	// Pretty output should have newlines
	if len(pretty) == 0 {
		t.Fatal("expected non-empty output")
	}
	hasNewline := false
	for _, b := range pretty {
		if b == '\n' {
			hasNewline = true
			break
		}
	}
	if !hasNewline {
		t.Error("expected pretty output to contain newlines")
	}

	compact, err := SerializeAvro(proj, false)
	if err != nil {
		t.Fatalf("serialize compact failed: %v", err)
	}
	// Compact should be shorter than pretty
	if len(compact) >= len(pretty) {
		t.Errorf("expected compact (%d) < pretty (%d)", len(compact), len(pretty))
	}
}

func TestExportAvroToFile_WritesValidJSON(t *testing.T) {
	input := []byte(`{"type":"enum","name":"Color","symbols":["RED","GREEN","BLUE"]}`)
	proj, _, err := LoadAvroFromBytes(input)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "color.avsc")

	if err := ExportAvroToFile(proj, outPath); err != nil {
		t.Fatalf("export failed: %v", err)
	}

	// Verify file exists and is valid JSON
	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("reading output failed: %v", err)
	}
	if len(content) == 0 {
		t.Fatal("expected non-empty file")
	}

	// Re-parse to verify valid JSON
	_, _, err = LoadAvroFromBytes(content)
	if err != nil {
		t.Fatalf("exported JSON is not valid avro: %v", err)
	}
}
