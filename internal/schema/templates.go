package schema

// TypeTemplates provides default Avro specs for creating new type nodes.
// Keys are type names; values are the raw JSON-equivalent Go structures.
var TypeTemplates = map[string]any{
	"union":  []any{"null"},
	"array":  map[string]any{"type": "array", "items": "string"},
	"map":    map[string]any{"type": "map", "values": "string"},
	"record": map[string]any{"type": "record", "name": "NewRecord", "namespace": "", "fields": []any{}},
	"enum":   map[string]any{"type": "enum", "name": "NewEnum", "symbols": []any{}},
	"fixed":  map[string]any{"type": "fixed", "name": "NewFixed", "size": float64(1)},
}

// NewFieldSpec returns a minimal field spec suitable for adding to a record.
func NewFieldSpec(name string, typeName string) map[string]any {
	return map[string]any{
		"name": name,
		"type": typeName,
	}
}
