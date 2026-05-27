package layoutcatalog

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrUnknownModule        = errors.New("layout module not found")
	ErrMissingRequiredField = errors.New("missing required field")
	ErrInvalidFieldValue    = errors.New("invalid field value")
)

func (c *Catalog) Render(name string, vars map[string]any) (string, error) {
	spec, ok := c.Get(name)
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnknownModule, name)
	}
	switch spec.BodyFormat {
	case BodyFormatRows:
		return renderRows(spec, vars)
	case BodyFormatJSONObject:
		return renderJSONFields(spec, vars, "object")
	case BodyFormatJSONArray:
		return renderJSONFields(spec, vars, "array")
	case BodyFormatFields, "":
		return renderFields(spec, vars)
	default:
		return "", fmt.Errorf("%w: unsupported body_format %q", ErrInvalidFieldValue, spec.BodyFormat)
	}
}

func renderFields(spec *LayoutSpec, vars map[string]any) (string, error) {
	var b strings.Builder
	fmt.Fprintf(&b, ":::%s\n", spec.Name)

	if spec.Fields != nil {
		for _, f := range spec.Fields.Required {
			val, ok := lookupString(vars, f.Name)
			if !ok || val == "" {
				return "", fmt.Errorf("%w: %s.%s", ErrMissingRequiredField, spec.Name, f.Name)
			}
			if err := checkEnum(f, val); err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "%s: %s\n", f.Name, val)
		}
		for _, f := range spec.Fields.Optional {
			val, ok := lookupString(vars, f.Name)
			if !ok || val == "" {
				continue
			}
			if err := checkEnum(f, val); err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "%s: %s\n", f.Name, val)
		}
	}
	b.WriteString(":::\n")
	return b.String(), nil
}

func renderJSONFields(spec *LayoutSpec, vars map[string]any, bodyKind string) (string, error) {
	obj := map[string]any{}
	if spec.Fields != nil {
		for _, f := range spec.Fields.Required {
			val, ok := lookupString(vars, f.Name)
			if !ok || val == "" {
				return "", fmt.Errorf("%w: %s.%s", ErrMissingRequiredField, spec.Name, f.Name)
			}
			if err := checkEnum(f, val); err != nil {
				return "", err
			}
			setJSONField(obj, f.Name, parseJSONFieldValue(val))
		}
		for _, f := range spec.Fields.Optional {
			val, ok := lookupString(vars, f.Name)
			if !ok || val == "" {
				continue
			}
			if err := checkEnum(f, val); err != nil {
				return "", err
			}
			setJSONField(obj, f.Name, parseJSONFieldValue(val))
		}
	}

	var body any = obj
	if bodyKind == "array" {
		body = []any{obj}
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	fmt.Fprintf(&b, ":::%s\n", spec.Name)
	fmt.Fprintf(&b, "%s\n", encoded)
	b.WriteString(":::\n")
	return b.String(), nil
}

func renderRows(spec *LayoutSpec, vars map[string]any) (string, error) {
	rowsRaw, ok := vars["rows"]
	if !ok {
		return "", fmt.Errorf("%w: %s.rows", ErrMissingRequiredField, spec.Name)
	}
	rows, ok := rowsRaw.([]any)
	if !ok || len(rows) == 0 {
		return "", fmt.Errorf("%w: %s.rows must be non-empty list", ErrInvalidFieldValue, spec.Name)
	}

	var b strings.Builder
	fmt.Fprintf(&b, ":::%s\n", spec.Name)

	if spec.Fields != nil {
		for _, f := range spec.Fields.Required {
			val, ok := lookupString(vars, f.Name)
			if !ok || val == "" {
				return "", fmt.Errorf("%w: %s.%s", ErrMissingRequiredField, spec.Name, f.Name)
			}
			if err := checkEnum(f, val); err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "%s: %s\n", f.Name, val)
		}
		for _, f := range spec.Fields.Optional {
			val, ok := lookupString(vars, f.Name)
			if !ok || val == "" {
				continue
			}
			if err := checkEnum(f, val); err != nil {
				return "", err
			}
			fmt.Fprintf(&b, "%s: %s\n", f.Name, val)
		}
	}

	delim := spec.Rows.Delimiter
	if delim == "" {
		delim = "|"
	}
	for i, row := range rows {
		cells, ok := row.([]any)
		if !ok {
			return "", fmt.Errorf("%w: %s.rows[%d] must be a list", ErrInvalidFieldValue, spec.Name, i)
		}
		if spec.Rows.MinColumns > 0 && len(cells) < spec.Rows.MinColumns {
			return "", fmt.Errorf("%w: %s.rows[%d] needs at least %d columns", ErrMissingRequiredField, spec.Name, i, spec.Rows.MinColumns)
		}
		strCells := make([]string, len(cells))
		for j, cell := range cells {
			strCells[j] = fmt.Sprintf("%v", cell)
		}
		fmt.Fprintln(&b, strings.Join(strCells, delim))
	}
	b.WriteString(":::\n")
	return b.String(), nil
}

func exampleJSONBodyKind(example string) string {
	body := strings.Split(example, "\n")
	for _, ln := range body {
		trimmed := strings.TrimSpace(ln)
		if trimmed == "" || strings.HasPrefix(trimmed, ":::") {
			continue
		}
		switch trimmed[0] {
		case '{':
			return "object"
		case '[':
			return "array"
		default:
			return ""
		}
	}
	return ""
}

func parseJSONFieldValue(val string) any {
	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return val
	}
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		var decoded any
		if err := json.Unmarshal([]byte(trimmed), &decoded); err == nil {
			return decoded
		}
	}
	return val
}

func setJSONField(obj map[string]any, key string, value any) {
	parts := strings.Split(key, ".")
	current := obj
	for _, part := range parts[:len(parts)-1] {
		next, ok := current[part].(map[string]any)
		if !ok {
			next = map[string]any{}
			current[part] = next
		}
		current = next
	}
	current[parts[len(parts)-1]] = value
}

func lookupString(vars map[string]any, key string) (string, bool) {
	v, ok := vars[key]
	if !ok {
		return "", false
	}
	switch t := v.(type) {
	case string:
		return t, true
	case fmt.Stringer:
		return t.String(), true
	default:
		return fmt.Sprintf("%v", v), true
	}
}

func checkEnum(f FieldSpec, val string) error {
	if len(f.Enum) == 0 {
		return nil
	}
	for _, allowed := range f.Enum {
		if allowed == val {
			return nil
		}
	}
	return fmt.Errorf("%w: %s must be one of %v, got %q", ErrInvalidFieldValue, f.Name, f.Enum, val)
}
