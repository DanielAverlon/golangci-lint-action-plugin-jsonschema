package rules

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var NoCommas = &analysis.Analyzer{
	Name: "nocommas",
	Doc:  "Disallow commas in JSON Schema fields",
	Run:  noCommasRun,
}

var allowedKeys = map[string]struct{}{
	"$id": {}, "$schema": {}, "$ref": {}, "$comment": {},
	"title": {}, "description": {}, "default": {}, "readOnly": {},
	"examples": {}, "multipleOf": {}, "maximum": {}, "exclusiveMaximum": {},
	"minimum": {}, "exclusiveMinimum": {}, "maxLength": {}, "minLength": {},
	"pattern": {}, "additionalItems": {}, "items": {}, "maxItems": {},
	"minItems": {}, "uniqueItems": {}, "contains": {}, "maxProperties": {},
	"minProperties": {}, "required": {}, "additionalProperties": {},
	"definitions": {}, "properties": {}, "patternProperties": {},
	"dependencies": {}, "propertyNames": {}, "const": {}, "enum": {},
	"type": {}, "format": {}, "contentMediaType": {}, "contentEncoding": {},
	"if": {}, "then": {}, "else": {}, "allOf": {}, "anyOf": {}, "oneOf": {}, "not": {},
}

func noCommasRun(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			field, ok := n.(*ast.Field)
			if !ok || field.Tag == nil {
				return true
			}

			// Example tag: `jsonschema:"title=Deprecation Date,description=The exact date, when deprecated."`
			tag := strings.Trim(field.Tag.Value, "`")
			jsonschemaTag := parseTagValue(tag, "jsonschema")
			if jsonschemaTag == "" {
				return true
			}

			// Check if the tags are allowed keys only
			for k, v := range parseKeyValuePairsLoose(jsonschemaTag) {
				if _, allowed := allowedKeys[k]; !allowed {
					pass.Report(analysis.Diagnostic{
						Pos:     field.Tag.Pos(),
						End:     field.Tag.End(),
						Message: "JSON Schema fields must use allowed keys",
					})
					return true
				}

				// Parse key=value pairs in the jsonschema tag
				if strings.Contains(v, ",") {
					noCommas := strings.ReplaceAll(v, ",", "\\,")
					newText := strings.ReplaceAll(tag, v, noCommas)
					pass.Report(analysis.Diagnostic{
						Pos:     field.Tag.Pos(),
						End:     field.Tag.End(),
						Message: "JSON Schema description fields should not contain commas",
						SuggestedFixes: []analysis.SuggestedFix{
							{
								Message: "Add escape for comma in description",
								TextEdits: []analysis.TextEdit{
									{
										Pos:     field.Tag.Pos() + 1,
										End:     field.Tag.End() - 1,
										NewText: []byte(newText),
									},
								},
							},
						},
					})

				}
			}

			return true
		})
	}
	return nil, nil
}

// parseTagValue extracts the quoted value for a tag key, e.g. jsonschema:"foo"
func parseTagValue(tag, key string) string {
	prefix := key + ":"
	start := strings.Index(tag, prefix)
	if start == -1 {
		return ""
	}
	start += len(prefix)
	if start >= len(tag) || tag[start] != '"' {
		return ""
	}
	end := strings.Index(tag[start+1:], `"`)
	if end == -1 {
		return ""
	}
	return tag[start+1 : start+1+end]
}

// parseKeyValuePairsLoose only splits on key=value boundaries
func parseKeyValuePairsLoose(s string) map[string]string {
	result := make(map[string]string)
	fields := strings.Split(s, ",")
	var currentKey string
	for _, part := range fields {
		//part = strings.TrimSpace(part)
		if eq := strings.Index(part, "="); eq != -1 {
			key := strings.TrimSpace(part[:eq])
			val := strings.TrimSpace(part[eq+1:])
			result[key] = val
			currentKey = key
		} else if currentKey != "" {
			// continuation of previous value (contains comma)
			result[currentKey] += "," + part
		}
	}
	return result
}
