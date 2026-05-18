package schema

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValidationResult represents the outcome of a validation check.
type ValidationResult struct {
	Valid   bool
	Message string
}

// OK returns a passing validation result.
func OK() ValidationResult {
	return ValidationResult{Valid: true}
}

// Fail returns a failing validation result with the given message.
func Fail(msg string) ValidationResult {
	return ValidationResult{Valid: false, Message: msg}
}

var avroNameRE = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// ValidateName checks an Avro name (field, record, enum, fixed).
func ValidateName(value string) ValidationResult {
	if strings.TrimSpace(value) == "" {
		return Fail("Name is required")
	}
	if !avroNameRE.MatchString(value) {
		return Fail("Name must start with [A-Za-z_] and contain only [A-Za-z0-9_]")
	}
	return OK()
}

// ValidateNamespace checks a dot-separated Avro namespace.
func ValidateNamespace(value string) ValidationResult {
	if strings.TrimSpace(value) == "" {
		return OK() // empty namespace is valid
	}
	segments := strings.Split(value, ".")
	for _, seg := range segments {
		if !avroNameRE.MatchString(seg) {
			return Fail("Each namespace segment must match [A-Za-z_][A-Za-z0-9_]*")
		}
	}
	return OK()
}

// ValidateSymbol checks an enum symbol name and ensures it's not a duplicate.
func ValidateSymbol(value string, existing []string) ValidationResult {
	nameResult := ValidateName(value)
	if !nameResult.Valid {
		return nameResult
	}
	for _, s := range existing {
		if s == value {
			return Fail("Duplicate symbol: " + value)
		}
	}
	return OK()
}

// ValidateFixedSize checks that a fixed size is a positive integer.
func ValidateFixedSize(value string) ValidationResult {
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return Fail("Fixed size must be an integer > 0")
	}
	return OK()
}

// ValidateDecimalPrecision checks that precision is a positive integer.
func ValidateDecimalPrecision(value string) ValidationResult {
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return Fail("Precision must be > 0")
	}
	return OK()
}

// ValidateDecimalScale checks that scale is non-negative and <= precision.
func ValidateDecimalScale(value string, precision string) ValidationResult {
	scale, err := strconv.Atoi(value)
	if err != nil || scale < 0 {
		return Fail("Scale must be >= 0")
	}
	prec, err := strconv.Atoi(precision)
	if err == nil && scale > prec {
		return Fail(fmt.Sprintf("Scale must be <= precision (%d)", prec))
	}
	return OK()
}

// ValidateUnionBranchAdd checks if a new type can be added to a union.
func ValidateUnionBranchAdd(newType string, existingBranches []string) ValidationResult {
	if newType == "union" {
		return Fail("Unions cannot be nested")
	}
	// Primitives and simple complex types cannot be duplicated in a union.
	noDupTypes := append(PrimitiveTypes, "array", "map")
	for _, t := range noDupTypes {
		if t == newType {
			for _, existing := range existingBranches {
				if existing == newType {
					return Fail("Union already contains type: " + newType)
				}
			}
			break
		}
	}
	return OK()
}
