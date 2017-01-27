package schema

type NameSchema struct {
	segments []string
}

func NewNameSchema(nameString string) NameSchema {
	var n NameSchema

	// Parse as necessary...

	return n
}

func (n NameSchema) Matches(other NameSchema) bool {
	return false
}

func (n NameSchema) IsPrefixOf(other NameSchema) bool {
	return false
}

type Schema struct {
	schema map[string]string // map from regex names to regex names
}
