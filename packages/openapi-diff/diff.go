package openapidiff

// Compare compares two API specifications and returns the
// detected changes, including breaking changes.
//
// An error is returned if a schema reference cannot be resolved.
func Compare(oldSpec, newSpec OpenAPISpec) (ComparisonResponse, error) {
	return CompareSpecsChecked(oldSpec, newSpec)
}
