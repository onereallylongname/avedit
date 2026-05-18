// Package schema defines Avro type system constants and templates.
package schema

// NodeKind identifies the structural role of a projection node.
type NodeKind string

const (
	KindSchema    NodeKind = "schema"
	KindRecord    NodeKind = "record"
	KindField     NodeKind = "field"
	KindPrimitive NodeKind = "primitive"
	KindNamed     NodeKind = "named"
	KindUnion     NodeKind = "union"
	KindArray     NodeKind = "array"
	KindMap       NodeKind = "map"
	KindEnum      NodeKind = "enum"
	KindFixed     NodeKind = "fixed"
)

// Slot identifies a structural position where children can be attached.
type Slot string

const (
	SlotSchemaRoot   Slot = "schemaRoot"
	SlotRecordFields Slot = "record.fields"
	SlotFieldType    Slot = "field.type"
	SlotArrayItems   Slot = "array.items"
	SlotMapValues    Slot = "map.values"
	SlotUnionBranch  Slot = "union.branch"
)

// PrimitiveTypes lists all Avro primitive type names.
var PrimitiveTypes = []string{
	"null", "boolean", "int", "long", "float", "double", "bytes", "string",
}

// ComplexTypes lists all Avro complex type names.
var ComplexTypes = []string{
	"record", "array", "map", "enum", "fixed", "union",
}

// LogicalTypes maps base types to their allowed logical types.
var LogicalTypes = map[string][]string{
	"int":    {"", "date", "time-millis"},
	"long":   {"", "timestamp-millis", "timestamp-micros", "timestamp-nanos", "local-timestamp-millis", "local-timestamp-micros", "local-timestamp-nanos", "time-micros"},
	"bytes":  {"", "decimal", "big-decimal"},
	"fixed":  {"", "decimal", "duration", "uuid"},
	"string": {"", "uuid"},
}

// LogicalTypeAttrs defines extra attributes needed for specific logical types.
var LogicalTypeAttrs = map[string]map[string]int{
	"decimal": {"scale": 0, "precision": 0},
}

// IsPrimitive returns true if the given type name is an Avro primitive.
func IsPrimitive(typeName string) bool {
	for _, p := range PrimitiveTypes {
		if p == typeName {
			return true
		}
	}
	return false
}

// IsComplex returns true if the given type name is an Avro complex type.
func IsComplex(typeName string) bool {
	for _, c := range ComplexTypes {
		if c == typeName {
			return true
		}
	}
	return false
}

// AllowedChildren defines which child kinds each parent kind can accept.
var AllowedChildren = map[NodeKind][]NodeKind{
	KindSchema: {KindRecord},
	KindRecord: {KindField},
	KindField:  {KindPrimitive, KindNamed, KindRecord, KindEnum, KindFixed, KindUnion, KindArray, KindMap},
	KindArray:  {KindPrimitive, KindNamed, KindRecord, KindEnum, KindFixed, KindUnion, KindArray, KindMap},
	KindMap:    {KindPrimitive, KindNamed, KindRecord, KindEnum, KindFixed, KindUnion, KindArray, KindMap},
	KindUnion:  {KindPrimitive, KindNamed, KindRecord, KindEnum, KindFixed, KindArray, KindMap},
}

// SlotAccepts defines which node roles are accepted in each slot.
var SlotAccepts = map[Slot][]NodeKind{
	SlotSchemaRoot:   {KindRecord, KindEnum, KindFixed, KindArray, KindMap, KindUnion, KindPrimitive, KindNamed},
	SlotRecordFields: {KindField},
	SlotFieldType:    {KindRecord, KindEnum, KindFixed, KindArray, KindMap, KindUnion, KindPrimitive, KindNamed},
	SlotArrayItems:   {KindRecord, KindEnum, KindFixed, KindArray, KindMap, KindUnion, KindPrimitive, KindNamed},
	SlotMapValues:    {KindRecord, KindEnum, KindFixed, KindArray, KindMap, KindUnion, KindPrimitive, KindNamed},
	SlotUnionBranch:  {KindRecord, KindEnum, KindFixed, KindArray, KindMap, KindPrimitive, KindNamed},
}

// SlotAcceptsKind checks if a slot accepts a given node kind.
func SlotAcceptsKind(slot Slot, kind NodeKind) bool {
	for _, k := range SlotAccepts[slot] {
		if k == kind {
			return true
		}
	}
	return false
}

// StandardNativeKeys lists the standard Avro attributes per node kind.
var StandardNativeKeys = map[NodeKind][]string{
	KindField:     {"name", "type", "doc", "default", "order", "aliases"},
	KindRecord:    {"type", "name", "namespace", "doc", "fields", "aliases"},
	KindEnum:      {"type", "name", "namespace", "doc", "symbols", "default", "aliases"},
	KindFixed:     {"type", "name", "namespace", "size", "doc", "aliases"},
	KindPrimitive: {"type", "logicalType", "precision", "scale"},
	KindNamed:     {},
	KindUnion:     {},
	KindArray:     {"type", "items", "default"},
	KindMap:       {"type", "values", "default"},
	KindSchema:    {"type", "name", "namespace", "doc", "fields", "aliases"},
}
