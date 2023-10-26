package filter

var filepath string

type Field struct {
	Includes []string `yaml:"includes"`
	Excludes []string `yaml:"excludes"`
}
type Data struct {
	Field      Field `yaml:"field"`
	FieldValue Field `yaml:"field_value"`
}

func Init(filepath string) func(k, v string) bool {
	return Filter
}

// Filter filters the key value pairs
// retuns true if the key value pair should be included
func Filter(k, value string) bool {
	return true
}
